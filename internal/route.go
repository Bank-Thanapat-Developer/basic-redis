package internal

import (
	"time"

	"github.com/Bank-Thanapat-Developer/basic-redis/internal/handlers"
	"github.com/Bank-Thanapat-Developer/basic-redis/pkg/cache"
	"github.com/Bank-Thanapat-Developer/basic-redis/pkg/middleware"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(
	app *fiber.App,
	cacheService cache.CacheService,
	jwtSecret string,
	authHandler *handlers.AuthHandler,
	itemHandler *handlers.ItemHandler,
	refItemTypeHandler *handlers.RefItemTypeHandler,
) {
	// Rate limit: 60 req/min ต่อ IP (ใช้กับทุก route)
	app.Use(middleware.RateLimit(cacheService, 60, time.Minute))

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OK"})
	})

	// ---- Auth (public) ----
	auth := v1.Group("/auth")
	auth.Post("/register", authHandler.Register)
	// Rate limit เข้มขึ้นสำหรับ login: 10 req/min เพื่อป้องกัน brute force
	auth.Post("/login", middleware.RateLimit(cacheService, 10, time.Minute), authHandler.Login)

	// ---- Protected routes (ต้องมี JWT) ----
	protected := v1.Group("", middleware.JWTAuth(jwtSecret))

	protected.Get("/me", authHandler.Me)

	items := protected.Group("/items")
	items.Post("/", itemHandler.CreateItem)
	items.Get("/list-with-redis", itemHandler.GetListItemsWithRedis)
	items.Get("/list-with-out-redis", itemHandler.GetListItemsWithOutRedis)
	items.Get("/:id", itemHandler.GetItemById)
	items.Put("/:id", itemHandler.UpdateItem)
	items.Delete("/:id", itemHandler.DeleteItem)

	refItemTypes := protected.Group("/ref-item-types")
	refItemTypes.Post("/", refItemTypeHandler.CreateRefItemType)
	refItemTypes.Get("/", refItemTypeHandler.GetListRefItemTypes)
	refItemTypes.Get("/:id", refItemTypeHandler.GetRefItemTypeById)
}
