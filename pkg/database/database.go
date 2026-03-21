package database

import (
	"fmt"
	"log"

	"github.com/YannisMuminov/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB(cfg *config.DatabaseConfig) error {
	dsn := cfg.DSN()

	var err error

	DB, err = sqlx.Connect("postgres", dsn)

	if err != nil {
		return fmt.Errorf("failed to connect database, %w", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	log.Println("Database connected successfully")

	return nil
}

func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
		log.Println("Database connection closed")
	}
}

func GetDB() *sqlx.DB {
	return DB
}
