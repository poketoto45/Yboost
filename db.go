package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TopGame struct {
	ID              uint   `gorm:"primaryKey"`
	SteamID         string `gorm:"index"` // Index pour des recherches rapides
	AppID           int
	Name            string
	PlaytimeForever int
}

var db *gorm.DB

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("La variable d'environnement DATABASE_URL est requise")
	}

	// Ajoute &TopGame{} ici
	db.AutoMigrate(&Note{}, &TopGame{})

	// AutoMigrate va créer ou mettre à jour la table "top_games"
	if err := db.AutoMigrate(&TopGame{}); err != nil {
	    log.Printf("Avertissement migration : %v", err)
	}


	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Impossible de se connecter à la base de données : %v", err)
	}


	// Récupère la connexion SQL sous-jacente
    sqlDB, err := db.DB()
    if err != nil {
        log.Printf("Erreur lors de la récupération de la connexion SQL : %v", err)
    } else {
        // Teste la connexion
        if err := sqlDB.Ping(); err != nil {
            log.Fatalf("Impossible de ping la base : %v", err)
        }

        // Infos sur la connexion
        log.Printf("✓ Connexion à la base établie")
        log.Printf("✓ Max open connections: %d", sqlDB.Stats().MaxOpenConnections)
        log.Printf("✓ Connexions actives: %d", sqlDB.Stats().OpenConnections)
    }

	// AutoMigrate crée la table si elle n'existe pas
	if err := db.AutoMigrate(&Note{}); err != nil {
		log.Printf("Avertissement migration : %v", err)
	}

	log.Println("Base de données connectée et migrée avec succès")
}
