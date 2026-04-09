package domains

import (
	"context"

	"github.com/Bank-Thanapat-Developer/basic-redis/internal/dto"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/entities"
)

type AuthUsecase interface {
	Register(ctx context.Context, req dto.RegisterRequest) error
	Login(ctx context.Context, req dto.LoginRequest) (*dto.TokenResponse, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	FindByUsername(ctx context.Context, username string) (*entities.User, error)
}
