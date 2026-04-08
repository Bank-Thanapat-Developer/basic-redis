package repositories

import (
	"context"
	"errors"

	"github.com/Bank-Thanapat-Developer/basic-redis/internal/domains"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/entities"
	"gorm.io/gorm"
)

type refItemTypeRepository struct {
	db *gorm.DB
}

func NewRefItemTypeRepository(db *gorm.DB) domains.RefItemTypeRepository {
	return &refItemTypeRepository{db: db}
}

func (r *refItemTypeRepository) Create(ctx context.Context, refItemType *entities.RefItemType) error {
	return r.db.WithContext(ctx).Create(refItemType).Error
}

func (r *refItemTypeRepository) GetList(ctx context.Context) ([]entities.RefItemType, error) {
	var refItemTypes []entities.RefItemType
	err := r.db.WithContext(ctx).Order("id ASC").Find(&refItemTypes).Error
	return refItemTypes, err
}

func (r *refItemTypeRepository) GetById(ctx context.Context, id string) (*entities.RefItemType, error) {
	var refItemType entities.RefItemType
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&refItemType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ref_item_type not found")
		}
		return nil, err
	}
	return &refItemType, nil
}
