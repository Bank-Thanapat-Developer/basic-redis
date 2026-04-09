package main

import (
	"log"

	"github.com/Bank-Thanapat-Developer/basic-redis/config"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/handlers"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/repositories"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/usecases"
	"github.com/Bank-Thanapat-Developer/basic-redis/pkg/cache"
	"github.com/Bank-Thanapat-Developer/basic-redis/pkg/database"
	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.NewPostgresDB(cfg.Postgres)
	if err != nil {
		log.Fatalf("failed to connect postgres: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	cacheService, err := cache.NewRedisCache(cfg.Redis)
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	// auth
	userRepo := repositories.NewUserRepository(db)
	authUsecase := usecases.NewAuthUsecase(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	authHandler := handlers.NewAuthHandler(authUsecase)

	// item
	itemRepo := repositories.NewItemRepository(db)
	itemUsecase := usecases.NewItemUsecase(itemRepo, cacheService)
	itemHandler := handlers.NewItemHandler(itemUsecase)

	// ref_item_type
	refItemTypeRepo := repositories.NewRefItemTypeRepository(db)
	refItemTypeUsecase := usecases.NewRefItemTypeUsecase(refItemTypeRepo)
	refItemTypeHandler := handlers.NewRefItemTypeHandler(refItemTypeUsecase)

	app := fiber.New(fiber.Config{
		ProxyHeader: fiber.HeaderXForwardedFor,
		TrustProxy:  true,
		TrustProxyConfig: fiber.TrustProxyConfig{
			Loopback: true,
			Private:  true,
		},
	})

	internal.SetupRoutes(app, cacheService, cfg.JWT.Secret, authHandler, itemHandler, refItemTypeHandler)

	log.Printf("server starting on port %s", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
