package domains

import (
	"context"

	"github.com/Bank-Thanapat-Developer/basic-redis/internal/dto"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/entities"
)

type ItemUsecase interface {
	Create(ctx context.Context, item dto.ItemCreateRequest) (int, error)
	GetItemById(ctx context.Context, id string) (*dto.ItemResponse, error)
	GetListItems(ctx context.Context, useRedis bool) ([]dto.ItemResponse, error)
	Update(ctx context.Context, id string, item dto.ItemUpdateRequest) (*dto.ItemResponse, error)
	Delete(ctx context.Context, id string) error
}

type ItemRepository interface {
	Create(ctx context.Context, item *entities.Item) error
	CheckDuplicateName(ctx context.Context, name string) (bool, error)

	GetItemById(ctx context.Context, id string) (*entities.ItemWithRefItemType, error)
	GetListItems(ctx context.Context) ([]entities.ItemWithRefItemType, error)
	CountListItems(ctx context.Context) (int64, error)

	Update(ctx context.Context, id string, item *entities.Item) error
	Delete(ctx context.Context, id string) error
}
