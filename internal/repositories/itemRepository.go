package repositories

import (
	"context"
	"errors"

	"github.com/Bank-Thanapat-Developer/basic-redis/internal/domains"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/entities"
	"gorm.io/gorm"
)

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) domains.ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) Create(ctx context.Context, item *entities.Item) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *itemRepository) CheckDuplicateName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Item{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

func (r *itemRepository) GetItemById(ctx context.Context, id string) (*entities.ItemWithRefItemType, error) {
	var item entities.Item
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item not found")
		}
		return nil, err
	}

	var refType entities.RefItemType
	r.db.WithContext(ctx).Where("id = ?", item.RefItemTypeID).First(&refType)

	return &entities.ItemWithRefItemType{Item: item, RefItemType: refType}, nil
}

func (r *itemRepository) GetListItems(ctx context.Context) ([]entities.ItemWithRefItemType, error) {
	var baseItems []entities.Item
	if err := r.db.WithContext(ctx).Where("is_active = true").Order("id ASC").Find(&baseItems).Error; err != nil {
		return nil, err
	}

	refTypeIDSet := make(map[int]struct{})
	for _, item := range baseItems {
		refTypeIDSet[item.RefItemTypeID] = struct{}{}
	}

	refTypeIDs := make([]int, 0, len(refTypeIDSet))
	for id := range refTypeIDSet {
		refTypeIDs = append(refTypeIDs, id)
	}

	refTypeMap := make(map[int]entities.RefItemType)
	if len(refTypeIDs) > 0 {
		var refTypes []entities.RefItemType
		if err := r.db.WithContext(ctx).Where("id IN ?", refTypeIDs).Find(&refTypes).Error; err != nil {
			return nil, err
		}
		for _, rt := range refTypes {
			refTypeMap[rt.ID] = rt
		}
	}

	result := make([]entities.ItemWithRefItemType, len(baseItems))
	for i, item := range baseItems {
		result[i] = entities.ItemWithRefItemType{
			Item:        item,
			RefItemType: refTypeMap[item.RefItemTypeID],
		}
	}
	return result, nil
}

func (r *itemRepository) CountListItems(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Item{}).Where("is_active = true").Count(&count).Error
	return count, err
}

func (r *itemRepository) Update(ctx context.Context, id string, item *entities.Item) error {
	result := r.db.WithContext(ctx).Model(&entities.Item{}).Where("id = ?", id).Updates(item)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("item not found")
	}
	return nil
}

func (r *itemRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Item{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("item not found")
	}
	return nil
}
