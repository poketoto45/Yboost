package main

import (
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL manquant")
	}

	// Supabase requiert SSL
	if !strings.Contains(dsn, "sslmode") {
		if strings.Contains(dsn, "?") {
			dsn += "&sslmode=require"
		} else {
			dsn += "?sslmode=require"
		}
	}

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Connexion DB impossible : %v", err)
	}

	if err := db.AutoMigrate(&TopGame{}); err != nil {
		log.Fatalf("Migration échouée : %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Erreur SQL DB : %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Ping DB échoué : %v", err)
	}

	log.Println("✓ Base de données connectée et migrée")
}
