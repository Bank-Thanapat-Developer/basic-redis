package handlers

import (
	"context"
	"errors"
	"time"

	"github.com/Bank-Thanapat-Developer/basic-redis/internal/domains"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/dto"
	"github.com/gofiber/fiber/v3"
)

type ItemHandler struct {
	itemUsecase domains.ItemUsecase
}

func NewItemHandler(itemUsecase domains.ItemUsecase) *ItemHandler {
	return &ItemHandler{itemUsecase: itemUsecase}
}

func (h *ItemHandler) CreateItem(c fiber.Ctx) error {
	var item dto.ItemCreateRequest

	if err := c.Bind().Body(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	id, err := h.itemUsecase.Create(ctx, item)
	if err != nil {
		if errors.Is(err, errors.New("name already exists")) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}

func (h *ItemHandler) GetItemById(c fiber.Ctx) error {
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	item, err := h.itemUsecase.GetItemById(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(item)
}

func (h *ItemHandler) GetListItemsWithRedis(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	items, err := h.itemUsecase.GetListItems(ctx, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(items)
}

func (h *ItemHandler) GetListItemsWithOutRedis(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	items, err := h.itemUsecase.GetListItems(ctx, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(items)
}
