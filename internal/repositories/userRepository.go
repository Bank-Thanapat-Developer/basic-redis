package repositories

import (
	"context"
	"errors"

	"github.com/Bank-Thanapat-Developer/basic-redis/internal/domains"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/entities"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domains.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // ไม่เจอ → return nil แทน error
		}
		return nil, err
	}
	return &user, nil
}
