package main

import (
	"sort"

	"gorm.io/gorm"
	"main.go/api"
)

type TopGame struct {
	ID              uint   `gorm:"primaryKey;autoIncrement;column:id"`
	SteamID         string `gorm:"column:steam_id;index;not null"`
	AppID           int    `gorm:"column:app_id;not null"`
	Name            string `gorm:"column:name;not null"`
	PlaytimeForever int    `gorm:"column:playtime_forever"`
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

func GetTopGamesFromDB(steamID string) ([]TopGame, error) {
	var games []TopGame
	result := db.Where("steam_id = ?", steamID).
		Order("playtime_forever DESC").
		Find(&games)
	return games, result.Error
}
