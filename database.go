package main

import (
	"sort"
	"main.go/api" // Vérifie que c'est bien ton chemin d'import
	"gorm.io/gorm"
)

// 1. Le Modèle (la structure de ta table)
type TopGame struct {
	ID              uint   `gorm:"primaryKey"`
	SteamID         string `gorm:"index"`
	AppID           int
	Name            string
	PlaytimeForever int
}

// 2. La fonction de synchronisation (le "C" et le "U" de ton CRUD)
func SyncTopGames(steamID string, allGames []api.OwnedGame) error {
	// TRI : On classe par temps de jeu décroissant
	sort.Slice(allGames, func(i, j int) bool {
		return allGames[i].PlaytimeForever > allGames[j].PlaytimeForever
	})

	// LIMITE : Top 5
	limit := 5
	if len(allGames) < 5 {
		limit = len(allGames)
	}
	top5 := allGames[:limit]

	// TRANSACTION : Pour éviter les données corrompues
	return db.Transaction(func(tx *gorm.DB) error {
		// On supprime les anciens jeux de cet utilisateur
		if err := tx.Where("steam_id = ?", steamID).Delete(&TopGame{}).Error; err != nil {
			return err
		}

		// On ajoute les nouveaux
		for _, g := range top5 {
			newEntry := TopGame{
				SteamID:         steamID,
				AppID:           g.AppID,
				Name:            g.Name,
				PlaytimeForever: g.PlaytimeForever,
			}
			if err := tx.Create(&newEntry).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
