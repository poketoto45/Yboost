package main

import (
	"sort"
	"main.go/api"
	"gorm.io/gorm"
)

type TopGame struct {
	ID              uint   `gorm:"primaryKey"`
	SteamID         string `gorm:"index"`
	AppID           int
	Name            string
	PlaytimeForever int
}

func SyncTopGames(steamID string, allGames []api.OwnedGame) error {
	sort.Slice(allGames, func(i, j int) bool {
		return allGames[i].PlaytimeForever > allGames[j].PlaytimeForever
	})

	limit := 5
	if len(allGames) < 5 {
		limit = len(allGames)
	}
	top5 := allGames[:limit]

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("steam_id = ?", steamID).Delete(&TopGame{}).Error; err != nil {
			return err
		}

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
