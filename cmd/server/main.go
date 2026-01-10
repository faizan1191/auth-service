package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/faizan1191/auth-service/internal/config"
	"github.com/faizan1191/auth-service/internal/router"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	r := router.SetupRouter()

	log.Println("server running on port", cfg.Port)
	r.Run(":" + cfg.Port)
}
