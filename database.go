package main

import (
	"sort"

	"gorm.io/gorm"
	"main.go/api"
)

type TopGame struct {
	ID              uint   `gorm:"primaryKey;autoIncrement"`
	SteamID         string `gorm:"index;not null"`
	AppID           int    `gorm:"not null"`
	Name            string `gorm:"not null"`
	PlaytimeForever int
}

func SyncTopGames(steamID string, allGames []api.OwnedGame) error {
	// Trier par temps de jeu décroissant
	sort.Slice(allGames, func(i, j int) bool {
		return allGames[i].PlaytimeForever > allGames[j].PlaytimeForever
	})

	limit := 5
	if len(allGames) < 5 {
		limit = len(allGames)
	}
	top5 := allGames[:limit]

	return db.Transaction(func(tx *gorm.DB) error {
		// DELETE : supprime les anciens enregistrements de ce joueur
		if err := tx.Where("steam_id = ?", steamID).Delete(&TopGame{}).Error; err != nil {
			return err
		}

		// CREATE : insère les nouveaux top 5
		for _, g := range top5 {
			entry := TopGame{
				SteamID:         steamID,
				AppID:           g.AppID,
				Name:            g.Name,
				PlaytimeForever: g.PlaytimeForever,
			}
			if err := tx.Create(&entry).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
