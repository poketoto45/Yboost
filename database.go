package main

import (
	"fmt"
	"sort"

	"main.go/api"
)

type SteamDB struct {
	ID      int64  `gorm:"primaryKey;column:id;autoIncrement"`
	SteamID string `gorm:"column:steam_id;uniqueIndex"`
	Game1   string `gorm:"column:game1"`
	Game2   string `gorm:"column:game2"`
	Game3   string `gorm:"column:game3"`
	Game4   string `gorm:"column:game4"`
	Game5   string `gorm:"column:game5"`
}

func (SteamDB) TableName() string { return "steamDB" }

func formatGame(g api.OwnedGame) string {
	return fmt.Sprintf("%s (%dh)", g.Name, g.PlaytimeForever/60)
}

func SyncTopGames(steamID string, allGames []api.OwnedGame) error {
	sort.Slice(allGames, func(i, j int) bool {
		return allGames[i].PlaytimeForever > allGames[j].PlaytimeForever
	})

	names := make([]string, 5)
	for i := 0; i < 5 && i < len(allGames); i++ {
		names[i] = formatGame(allGames[i])
	}

	// Vérifie si une ligne existe déjà pour ce steamID
	var count int64
	db.Model(&SteamDB{}).Where("steam_id = ?", steamID).Count(&count)

	if count == 0 {
		// INSERT
		entry := SteamDB{
			SteamID: steamID,
			Game1:   names[0],
			Game2:   names[1],
			Game3:   names[2],
			Game4:   names[3],
			Game5:   names[4],
		}
		return db.Create(&entry).Error
	}

	// UPDATE direct sans passer par Save
	return db.Model(&SteamDB{}).
		Where("steam_id = ?", steamID).
		Updates(map[string]interface{}{
			"game1": names[0],
			"game2": names[1],
			"game3": names[2],
			"game4": names[3],
			"game5": names[4],
		}).Error
}

func GetTopGamesFromDB(steamID string) (*SteamDB, error) {
	var row SteamDB
	err := db.Where("steam_id = ?", steamID).First(&row).Error
	return &row, err
}

func DeleteTopGames(steamID string) error {
	return db.Where("steam_id = ?", steamID).Delete(&SteamDB{}).Error
}
