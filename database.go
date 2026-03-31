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

	// UPDATE & CREATE (Upsert)
	// On cherche par SteamID, si trouvé on Update, sinon on Create
	var entry SteamDB
	result := db.Where("steam_id = ?", steamID).First(&entry)

	entry.SteamID = steamID
	entry.Game1, entry.Game2, entry.Game3, entry.Game4, entry.Game5 = names[0], names[1], names[2], names[3], names[4]

	if result.Error != nil {
		return db.Create(&entry).Error
	}
	return db.Save(&entry).Error
}

func GetTopGamesFromDB(steamID string) (*SteamDB, error) {
	var row SteamDB
	err := db.Where("steam_id = ?", steamID).First(&row).Error
	return &row, err
}

func DeleteTopGames(steamID string) error {
	return db.Where("steam_id = ?", steamID).Delete(&SteamDB{}).Error
}
