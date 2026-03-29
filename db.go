package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Le struct TopGame est supprimé d'ici car il est déjà dans database.go
var db *gorm.DB

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("La variable d'environnement DATABASE_URL est requise")
	}

	var err error
	// CHANGEMENT : On ouvre la connexion AVANT de migrer
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Impossible de se connecter à la base de données : %v", err)
	}

	// CHANGEMENT : Supprimé &Note{} qui n'existe pas, gardé uniquement &TopGame{}
	if err := db.AutoMigrate(&TopGame{}); err != nil {
		log.Printf("Avertissement migration : %v", err)
	}

	// Récupère la connexion SQL sous-jacente
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Erreur lors de la récupération de la connexion SQL : %v", err)
	} else {
		if err := sqlDB.Ping(); err != nil {
			log.Fatalf("Impossible de ping la base : %v", err)
		}
		log.Printf("✓ Connexion à la base établie")
		log.Printf("✓ Max open connections: %d", sqlDB.Stats().MaxOpenConnections)
		log.Printf("✓ Connexions actives: %d", sqlDB.Stats().OpenConnections)
	}

	log.Println("Base de données connectée et migrée avec succès")
}
