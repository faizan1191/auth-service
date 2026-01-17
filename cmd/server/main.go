package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"

	"github.com/faizan1191/auth-service/internal/auth"
	"github.com/faizan1191/auth-service/internal/config"
	"github.com/faizan1191/auth-service/internal/db"
	"github.com/faizan1191/auth-service/internal/redis"
	"github.com/faizan1191/auth-service/internal/router"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	// connect to PostgreSQL
	database, err := db.NewPostgres(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	// init repository
	userRepo := auth.NewRepository(database)

	// init redis client
	redisClient := redis.NewRedisClient()

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	// init handler with repo
	authHandler := auth.NewHandler(userRepo, redisClient)

	// setup router with handlers
	r := router.SetupRouter(authHandler)

	log.Println("server running on port", cfg.Port)

	r.Run(":" + cfg.Port)
}
