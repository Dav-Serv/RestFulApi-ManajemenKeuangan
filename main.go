package main

import (
	"log"

	"RestFulApi-ManajemenKeuangan/config"
	"RestFulApi-ManajemenKeuangan/internal/database"
	"RestFulApi-ManajemenKeuangan/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	gin.SetMode(cfg.GinMode)

	database.Connect(cfg.DBPath)

	router := routes.Setup()

	log.Printf("[server] berjalan di http://localhost:%s\n", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("[server] gagal start: %v", err)
	}
}
