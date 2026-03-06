package client

import (
	"context"
	"errors"

	"github.com/brucechen520/kuji-go/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
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

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	// 使用 .Where 加上條件，並透過 .First 抓取單一筆
	// 記得傳入 ctx 以利於追蹤與超時控制
	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 沒找到時回傳 nil, nil 是一種常見做法
		}
		return nil, err
	}
	return &user, nil
}
