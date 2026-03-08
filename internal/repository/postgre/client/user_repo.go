package client

import (
	"errors"

	"github.com/brucechen520/kuji-go/internal/model"
	"github.com/brucechen520/kuji-go/internal/pkg/core"
	"gorm.io/gorm"
)

var _ UserRepository = (*userRepository)(nil)

type UserRepository interface {
	GetByEmail(ctx core.Context, email string) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepo 回傳的是介面，這是解耦的關鍵
func NewUserRepo(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetByEmail(ctx core.Context, email string) (*model.User, error) {
	var user model.User
	// 使用 .Where 加上條件，並透過 .First 抓取單一筆
	// 記得傳入 ctx 以利於追蹤與超時控制
	err := r.db.WithContext(ctx.StdContext()).Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 沒找到時回傳 nil, nil 是一種常見做法
		}
		return nil, err
	}
	return &user, nil
}
