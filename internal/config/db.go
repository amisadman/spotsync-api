package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg *Config) *gorm.DB {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: cfg.Dsn,
	}),&gorm.Config{
		TranslateError: true,
	})
	

	if err!= nil{
		log.Fatal("Failed to connect database")

	}

	println("Database connection successfull")
	return db
}