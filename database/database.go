package database

import (
	"fmt"
	"log"
	"todo-api/config"
	"todo-api/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(cfg *config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Jakarta",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("gagal terhubung ke gudang utama:", err)
	}

	err = db.AutoMigrate(&domain.Todo{})
	if err != nil {
		log.Fatal("gagal migrasi struktur gedung:", err)
	}

	DB = db
	log.Println("koneksi ke gudang utama berhasil!")

}
