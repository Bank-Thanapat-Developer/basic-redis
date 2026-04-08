package internal

import (
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/handlers"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App, itemHandler *handlers.ItemHandler, refItemTypeHandler *handlers.RefItemTypeHandler) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OK"})
	})

	items := v1.Group("/items")
	items.Post("/", itemHandler.CreateItem)
	items.Get("/list-with-redis", itemHandler.GetListItemsWithRedis)
	items.Get("/list-with-out-redis", itemHandler.GetListItemsWithOutRedis)
	items.Get("/:id", itemHandler.GetItemById)
	items.Put("/:id", itemHandler.UpdateItem)
	items.Delete("/:id", itemHandler.DeleteItem)

	refItemTypes := v1.Group("/ref-item-types")
	refItemTypes.Post("/", refItemTypeHandler.CreateRefItemType)
	refItemTypes.Get("/", refItemTypeHandler.GetListRefItemTypes)
	refItemTypes.Get("/:id", refItemTypeHandler.GetRefItemTypeById)
}
