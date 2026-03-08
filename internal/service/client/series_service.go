package client

import (
	"fmt"
	"strconv"

	"go.uber.org/zap"

	"github.com/brucechen520/kuji-go/internal/config"
	"github.com/brucechen520/kuji-go/internal/dto"
	"github.com/brucechen520/kuji-go/internal/dto/mapper"
	"github.com/brucechen520/kuji-go/internal/model"
	"github.com/brucechen520/kuji-go/internal/pkg/core"
	"github.com/brucechen520/kuji-go/internal/repository/postgre/client"
	"github.com/brucechen520/kuji-go/internal/repository/redis"
	"golang.org/x/sync/singleflight"
)

type SeriesService struct {
	seriesRepo client.SeriesRepository
	kujiStore  redis.KujiStore
	sf         singleflight.Group
}

func NewSeriesService(series client.SeriesRepository, kuji redis.KujiStore, cfg *config.AuthConfig) *SeriesService {
	return &SeriesService{seriesRepo: series, kujiStore: kuji}
}

func (s *SeriesService) GetSeriesById(ctx core.Context, id uint) (*dto.SeriesDetailDTO, error) {
	// 1. 先撈靜態 Meta (系列/箱子/獎項名稱)
	// 這裡我們用 singleflight 防止 Cache 擊穿
	series, err, _ := s.sf.Do("series_meta:"+strconv.Itoa(int(id)), func() (interface{}, error) {
		// 嘗試從 Redis 撈完整結構
		data, err := s.kujiStore.GetSeriesMeta(ctx, id)
		if err == nil {
			ctx.GetLogger().Info("[Cache Hit] Loaded from Redis", zap.Uint("series_id", id))
			return data, nil
		}

		ctx.GetLogger().Info("[Cache Miss] Fetching from DB", zap.Uint("series_id", id))

		// 如果 Redis 沒有，去 DB 撈
		dbData, err := s.seriesRepo.GetSeriesById(ctx, id)
		// 系統層錯誤：記錄日誌並向上拋出
		if err != nil {
			ctx.GetLogger().Error("DB Error", zap.Error(err))
			return nil, err
		}

		// 業務層錯誤：資料不存在
		if dbData == nil {
			return nil, fmt.Errorf("series %d not found", id)
		}

		// 回寫 Redis
		s.kujiStore.SetSeriesMeta(ctx, id, dbData)
		return dbData, nil
	})
	if err != nil {
		ctx.GetLogger().Error("GetSeriesMeta Error", zap.Error(err))
		return nil, err
	}

	// 2. 轉換成 API 用的 DTO (過濾掉所有敏感與無關欄位)
	// 此時前端只會看到你在 DTO 定義的那幾個 json 欄位
	dtoSeries := mapper.MapSeriesToDetailDTO(series.(*model.Series))

	// 2. 組合動態庫存
	for i := range dtoSeries.Boxes {
		// 從 Redis 拿 map[string]int
		inv, err := s.kujiStore.GetBoxInventory(ctx, dtoSeries.Boxes[i].ID)

		// 如果 Redis 沒資料 (inv 為 nil)
		if err == nil && inv == nil {
			// 去 DB 補貨
			prizes, dbErr := s.seriesRepo.GetBoxInventoryById(ctx, dtoSeries.Boxes[i].ID)
			if dbErr != nil { // 1. 處理技術錯誤
				ctx.GetLogger().Error("DB Error", zap.Error(err))
				return nil, err
			}

			// 業務層錯誤：資料不存在
			if len(prizes) == 0 {
				return nil, fmt.Errorf("Prizes %d not found", id)
			}

			if dbErr == nil && len(prizes) > 0 {
				// 轉換：將 DB 回傳的 []model.Prize 轉成 map 寫回 Redis
				inv = make(map[string]int, len(prizes))
				for _, p := range prizes {
					inv[strconv.Itoa(int(p.ID))] = p.RemainingQuantity
				}
				err = s.kujiStore.SetBoxInventory(ctx, dtoSeries.Boxes[i].ID, inv)
				if err != nil {
					ctx.GetLogger().Error("SetBoxInventory Error", zap.Error(err))
				}
			} else {
				inv = make(map[string]int) // 補貨失敗，設為空 map
			}
		}

		// 將 Hash 的值更新到 Series 結構中
		for j := range dtoSeries.Boxes[i].Prizes {
			prizeID := strconv.Itoa(int(dtoSeries.Boxes[i].Prizes[j].ID))
			if count, ok := inv[prizeID]; ok {
				dtoSeries.Boxes[i].Prizes[j].RemainingQuantity = count
			} else {
				dtoSeries.Boxes[i].Prizes[j].RemainingQuantity = 0
			}
		}
	}

	return dtoSeries, nil
}
