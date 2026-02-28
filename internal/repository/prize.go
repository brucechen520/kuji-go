package repository

import (
	"errors"

	"github.com/brucechen520/kuji-go/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// PrizeRepo 介面定義
type PrizeRepo interface {
	GetByID(id uint) (*models.Prize, error)
	CreatePrize(prize *models.Prize) error
	UpdateStock(id uint, num int) error
}

// 實作物件
type prizeRepo struct {
	db  *gorm.DB // 這裡持有 GORM 的連線
	rdb *redis.Client
}

// NewPrizeRepo 初始化
func NewPrizeRepo(db *gorm.DB, rdb *redis.Client) PrizeRepo {
	return &prizeRepo{db: db, rdb: rdb}
}

// GetByID 查詢單一獎品
func (r *prizeRepo) GetByID(id uint) (*models.Prize, error) {
	var prize models.Prize
	// First 會根據主鍵查詢，如果找不到會回傳 gorm.ErrRecordNotFound
	if err := r.db.First(&prize, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 或者回傳自定義的 NotFound Error
		}
		return nil, err
	}
	return &prize, nil
}

// CreatePrize 建立新獎品
func (r *prizeRepo) CreatePrize(prize *models.Prize) error {
	return r.db.Create(prize).Error
}

// UpdateStock 更新庫存 (建議使用原子操作)
func (r *prizeRepo) UpdateStock(id uint, num int) error {
	// 使用 Model 指定表，Where 指定條件，Update 執行變更
	// 這裡示範基本的 GORM 更新；若需併發安全，建議改用 gorm.Expr("stock + ?", num)
	return r.db.Model(&models.Prize{}).
		Where("id = ?", id).
		Update("stock", gorm.Expr("stock + ?", num)).
		Error
}
