package main

import (
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL manquant")
	}

	if !strings.Contains(dsn, "sslmode") {
		if strings.Contains(dsn, "?") {
			dsn += "&sslmode=require"
		} else {
			dsn += "?sslmode=require"
		}
	}

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Connexion DB impossible : %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Erreur SQL DB : %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Ping DB échoué : %v", err)
	}

	log.Println("✓ Base de données connectée")
}
