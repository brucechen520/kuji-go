package client

import (
	"errors"

	"github.com/brucechen520/kuji-go/internal/model"
	"github.com/brucechen520/kuji-go/internal/pkg/core"
	"gorm.io/gorm"
)

var _ SeriesRepository = (*seriesRepository)(nil)

type SeriesRepository interface {
	GetSeriesById(ctx core.Context, id uint) (*model.Series, error)
	GetBoxInventoryById(ctx core.Context, id uint) ([]model.Prize, error)
}

type seriesRepository struct {
	db *gorm.DB
}

// NewUserRepo 回傳的是介面，這是解耦的關鍵
func NewSeriesRepo(db *gorm.DB) SeriesRepository {
	return &seriesRepository{
		db: db,
	}
}

func (u *seriesRepository) GetSeriesById(ctx core.Context, id uint) (*model.Series, error) {
	var series model.Series
	// 使用 .Where 加上條件，並透過 .First 抓取單一筆
	// 記得傳入 ctx.StdContext() 以利於追蹤與超時控制
	err := u.db.WithContext(ctx.StdContext()).
		Select("id, name, price, description").
		Preload("Boxes", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, series_id, location_name")
		}).
		Preload("Boxes.Prizes", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, box_id, level, name, initial_quantity") // 注意：這裡不取 RemainingQuantity
		}).
		First(&series, id).Error

	// 判斷邏輯：
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // 資料庫回傳空，對外宣告「資料不存在」但「無系統錯誤」
	}
	if err != nil {
		return nil, err // 真正系統錯誤 (連線中斷、SQL 語法錯)
	}

	return &series, nil
}

func (u *seriesRepository) GetBoxInventoryById(ctx core.Context, id uint) ([]model.Prize, error) {
	var box model.Box
	err := u.db.WithContext(ctx.StdContext()).
		Select("id").
		Where("id = ? AND remain_quantity > 0", id).
		Preload("Prizes", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, box_id, level, name, remaining_quantity") // box_id 必須＋才能正確 join
		}).
		First(&box).Error

	// 判斷邏輯：
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // 資料庫回傳空，對外宣告「資料不存在」但「無系統錯誤」
	}

	if err != nil {
		return nil, err // 真正系統錯誤 (連線中斷、SQL 語法錯)
	}

	return box.Prizes, nil
}
