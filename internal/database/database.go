package database

import (
	"log"

	"RestFulApi-ManajemenKeuangan/internal/models"

	"github.com/glebarez/sqlite" // driver sqlite pure-Go (tanpa CGO)
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect membuka koneksi ke file SQLite dan menjalankan AutoMigrate.
func Connect(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		log.Fatalf("[database] gagal konek ke sqlite: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Transaction{}); err != nil {
		log.Fatalf("[database] gagal migrate: %v", err)
	}

	log.Println("[database] terkoneksi & migrasi selesai:", dbPath)
	DB = db
	return db
}