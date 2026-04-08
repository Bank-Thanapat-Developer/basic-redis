package usecases

import (
	"context"
	"time"

	"github.com/Bank-Thanapat-Developer/basic-redis/internal/domains"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/dto"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/entities"
)

type refItemTypeUsecase struct {
	refItemTypeRepository domains.RefItemTypeRepository
}

func NewRefItemTypeUsecase(repo domains.RefItemTypeRepository) domains.RefItemTypeUsecase {
	return &refItemTypeUsecase{refItemTypeRepository: repo}
}

func (u *refItemTypeUsecase) Create(ctx context.Context, req dto.RefItemTypeCreateRequest) (int, error) {
	data := &entities.RefItemType{
		Name:      req.Name,
		CreatedAt: time.Now(),
	}
	if err := u.refItemTypeRepository.Create(ctx, data); err != nil {
		return 0, err
	}
	return data.ID, nil
}

func (u *refItemTypeUsecase) GetList(ctx context.Context) ([]dto.RefItemTypeResponse, error) {
	refTypes, err := u.refItemTypeRepository.GetList(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]dto.RefItemTypeResponse, len(refTypes))
	for i, rt := range refTypes {
		response[i] = dto.RefItemTypeResponse{
			ID:   rt.ID,
			Name: rt.Name,
		}
	}
	return response, nil
}

func (u *refItemTypeUsecase) GetById(ctx context.Context, id string) (*dto.RefItemTypeResponse, error) {
	rt, err := u.refItemTypeRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.RefItemTypeResponse{
		ID:   rt.ID,
		Name: rt.Name,
	}, nil
}
