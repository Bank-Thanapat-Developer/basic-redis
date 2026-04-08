package domains

import (
	"context"

	"github.com/Bank-Thanapat-Developer/basic-redis/internal/dto"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/entities"
)

type RefItemTypeUsecase interface {
	Create(ctx context.Context, req dto.RefItemTypeCreateRequest) (int, error)
	GetList(ctx context.Context) ([]dto.RefItemTypeResponse, error)
	GetById(ctx context.Context, id string) (*dto.RefItemTypeResponse, error)
}

type RefItemTypeRepository interface {
	Create(ctx context.Context, refItemType *entities.RefItemType) error
	GetList(ctx context.Context) ([]entities.RefItemType, error)
	GetById(ctx context.Context, id string) (*entities.RefItemType, error)
}
