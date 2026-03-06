// internal/repository/redis/kuji_repo.go
package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/brucechen520/kuji-go/internal/model"
	"github.com/brucechen520/kuji-go/internal/pkg"
	"github.com/redis/go-redis/v9"
)

type KujiStore interface {
	GetBoxInventory(ctx context.Context, boxID uint) (map[string]int, error)
	SetBoxInventory(ctx context.Context, boxID uint, inv map[string]int) error
	GetSeriesMeta(ctx context.Context, seriesID uint) (*model.Series, error)
	SetSeriesMeta(ctx context.Context, seriesID uint, series *model.Series) error
}

type kujiStore struct {
	client *redis.Client
}

func NewKujiStore(client *redis.Client) KujiStore {
	return &kujiStore{client: client}
}

func (r *kujiStore) GetBoxInventory(ctx context.Context, boxID uint) (map[string]int, error) {
	key := pkg.GetBoxInventoryKey(boxID)
	// HGetAll 一次拿走這個箱子所有獎項的剩餘數
	results, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err // 真正的網路錯誤或 Redis 異常
	}

	// 關鍵判斷：如果長度為 0，代表 Redis 沒有該箱子的庫存資料
	if len(results) == 0 {
		return nil, nil // 告知 Service 去 DB 補貨
	}

	// 將 string 的 map 轉換成 int
	inv := make(map[string]int, len(results))
	for k, v := range results {
		// 這裡假設 Field 是 prizeID，Value 是 count
		inv[k] = pkg.ParseInt(v)
	}
	return inv, nil
}

func (r *kujiStore) SetBoxInventory(ctx context.Context, boxID uint, inv map[string]int) error {
	key := pkg.GetBoxInventoryKey(boxID)

	// 1. 建立 map[string]interface{}，這符合 go-redis 的 HSet 需求
	values := make(map[string]interface{}, len(inv))
	for prizeID, count := range inv {
		// 直接存入 interface{}，go-redis 會自動幫你轉成 string 寫入 Redis Hash
		values[prizeID] = count
	}

	// HSet 支援一次傳入一個 map，自動將所有 Field-Value 對寫入
	return r.client.HSet(ctx, key, values).Err()
}

func (r *kujiStore) GetSeriesMeta(ctx context.Context, seriesID uint) (*model.Series, error) {
	key := pkg.GetSeriesMetaKey(seriesID)

	// 1. 從 Redis 讀取 JSON 字串
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err // 如果沒資料或連線失敗，回傳錯誤讓 Service 去 DB 撈
	}

	// 2. 反序列化成結構體
	var series model.Series
	if err := json.Unmarshal(val, &series); err != nil {
		return nil, err
	}
	return &series, nil
}

func (r *kujiStore) SetSeriesMeta(ctx context.Context, seriesID uint, series *model.Series) error {
	key := pkg.GetSeriesMetaKey(seriesID)

	// 1. 序列化
	data, err := json.Marshal(series)
	if err != nil {
		return err
	}

	// 2. 寫入 Redis (設定 TTL，例如 1 小時過期)
	return r.client.Set(ctx, key, data, 1*time.Hour).Err()
}
