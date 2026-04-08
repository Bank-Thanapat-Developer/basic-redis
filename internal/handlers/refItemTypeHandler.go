package handlers

import (
	"context"
	"time"

	"github.com/Bank-Thanapat-Developer/basic-redis/internal/domains"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/dto"
	"github.com/gofiber/fiber/v3"
)

type RefItemTypeHandler struct {
	refItemTypeUsecase domains.RefItemTypeUsecase
}

func NewRefItemTypeHandler(refItemTypeUsecase domains.RefItemTypeUsecase) *RefItemTypeHandler {
	return &RefItemTypeHandler{refItemTypeUsecase: refItemTypeUsecase}
}

func (h *RefItemTypeHandler) CreateRefItemType(c fiber.Ctx) error {
	var req dto.RefItemTypeCreateRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	id, err := h.refItemTypeUsecase.Create(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}

func (h *RefItemTypeHandler) GetListRefItemTypes(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	refTypes, err := h.refItemTypeUsecase.GetList(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(refTypes)
}

func (h *RefItemTypeHandler) GetRefItemTypeById(c fiber.Ctx) error {
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	refType, err := h.refItemTypeUsecase.GetById(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(refType)
}
