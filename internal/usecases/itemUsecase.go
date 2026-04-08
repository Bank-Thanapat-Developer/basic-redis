package usecases

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Bank-Thanapat-Developer/basic-redis/internal/domains"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/dto"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/entities"
	"github.com/Bank-Thanapat-Developer/basic-redis/pkg/cache"
)

type itemUsecase struct {
	itemRepository domains.ItemRepository
	cache          cache.CacheService
}

func NewItemUsecase(itemRepository domains.ItemRepository, cache cache.CacheService) domains.ItemUsecase {
	return &itemUsecase{itemRepository: itemRepository, cache: cache}
}

func (u *itemUsecase) Create(ctx context.Context, item dto.ItemCreateRequest) (int, error) {
	isDuplicate, err := u.itemRepository.CheckDuplicateName(ctx, item.Name)
	if err != nil {
		return 0, err
	}
	if isDuplicate {
		return 0, errors.New("name already exists")
	}

	itemData := entities.Item{
		Name:          item.Name,
		Price:         item.Price,
		IsActive:      item.IsActive,
		RefItemTypeID: item.RefItemTypeID,
		CreatedAt:     time.Now(),
	}

	err = u.itemRepository.Create(ctx, &itemData)
	if err != nil {
		return 0, err
	}

	// invalidate list cache เพราะมีข้อมูลใหม่เข้ามา
	if err := u.cache.Delete(ctx, "items:list"); err != nil {
		log.Println("failed to invalidate items:list cache:", err)
	}

	return itemData.ID, nil
}

func (u *itemUsecase) GetItemById(ctx context.Context, id string) (*dto.ItemResponse, error) {
	cacheKey := fmt.Sprintf("item:%s", id)

	var cached dto.ItemResponse
	if err := u.cache.GetObject(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	item, err := u.itemRepository.GetItemById(ctx, id)
	if err != nil {
		return nil, err
	}

	response := dto.ItemResponse{
		ID:       item.ID,
		Name:     item.Name,
		Price:    item.Price,
		IsActive: item.IsActive,
		RefItemType: &dto.RefItemTypeResponse{
			ID:   item.RefItemType.ID,
			Name: item.RefItemType.Name,
		},
	}

	if err := u.cache.Set(ctx, cacheKey, response, 10*time.Minute); err != nil {
		log.Println("failed to set cache:", err)
	}

	return &response, nil
}

func (u *itemUsecase) GetListItems(ctx context.Context, useRedis bool) ([]dto.ItemResponse, error) {

	cacheKey := "items:list"

	if useRedis {
		var cacheValue []dto.ItemResponse
		if err := u.cache.GetObject(ctx, cacheKey, &cacheValue); err == nil && len(cacheValue) > 0 {
			return cacheValue, nil
		}
	}

	items, err := u.itemRepository.GetListItems(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]dto.ItemResponse, len(items))
	for i, item := range items {
		response[i] = dto.ItemResponse{
			ID:       item.ID,
			Name:     item.Name,
			Price:    item.Price,
			IsActive: item.IsActive,
			RefItemType: &dto.RefItemTypeResponse{
				ID:   item.RefItemType.ID,
				Name: item.RefItemType.Name,
			},
		}
	}

	if useRedis {
		ttl := 5 * time.Minute
		err = u.cache.Set(ctx, cacheKey, response, ttl)
		if err != nil {
			log.Println("failed to set cache:", err)
		}
	}

	return response, nil
}

func (u *itemUsecase) Update(ctx context.Context, id string, req dto.ItemUpdateRequest) (*dto.ItemResponse, error) {
	itemData := &entities.Item{
		Name:          req.Name,
		Price:         req.Price,
		IsActive:      req.IsActive,
		RefItemTypeID: req.RefItemTypeID,
	}

	if err := u.itemRepository.Update(ctx, id, itemData); err != nil {
		return nil, err
	}

	// invalidate cache ของ item นี้ และ list cache
	cacheKey := fmt.Sprintf("item:%s", id)
	if err := u.cache.Delete(ctx, cacheKey); err != nil {
		log.Println("failed to invalidate item cache:", err)
	}
	if err := u.cache.Delete(ctx, "items:list"); err != nil {
		log.Println("failed to invalidate items:list cache:", err)
	}

	// ดึงข้อมูลล่าสุดกลับไป (cache จะถูก set ใหม่ใน GetItemById)
	return u.GetItemById(ctx, id)
}

func (u *itemUsecase) Delete(ctx context.Context, id string) error {
	if err := u.itemRepository.Delete(ctx, id); err != nil {
		return err
	}

	// invalidate cache ของ item นี้ และ list cache
	cacheKey := fmt.Sprintf("item:%s", id)
	if err := u.cache.Delete(ctx, cacheKey); err != nil {
		log.Println("failed to invalidate item cache:", err)
	}
	if err := u.cache.Delete(ctx, "items:list"); err != nil {
		log.Println("failed to invalidate items:list cache:", err)
	}

	return nil
}
