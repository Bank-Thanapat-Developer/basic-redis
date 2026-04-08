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
	// load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// initialize postgres database
	db, err := database.NewPostgresDB(cfg.Postgres)
	if err != nil {
		log.Fatalf("failed to connect postgres: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	// initialize redis client
	cacheService, err := cache.NewRedisCache(cfg.Redis)
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	// item
	itemRepo := repositories.NewItemRepository(db)
	itemUsecase := usecases.NewItemUsecase(itemRepo, cacheService)
	itemHandler := handlers.NewItemHandler(itemUsecase)

	// ref_item_type
	refItemTypeRepo := repositories.NewRefItemTypeRepository(db)
	refItemTypeUsecase := usecases.NewRefItemTypeUsecase(refItemTypeRepo)
	refItemTypeHandler := handlers.NewRefItemTypeHandler(refItemTypeUsecase)

	app := fiber.New()

	internal.SetupRoutes(app, itemHandler, refItemTypeHandler)

	log.Printf("server starting on port %s", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
