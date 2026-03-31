package main

import (
	"fmt"
	"sort"

	"main.go/api"
)

type SteamDB struct {
	ID      int64  `gorm:"primaryKey;column:id;autoIncrement"`
	SteamID string `gorm:"column:steam_id"`
	Game1   string `gorm:"column:game1"`
	Game2   string `gorm:"column:game2"`
	Game3   string `gorm:"column:game3"`
	Game4   string `gorm:"column:game4"`
	Game5   string `gorm:"column:game5"`
}

func (SteamDB) TableName() string {
	return "steamDB"
}

func formatGame(g api.OwnedGame) string {
	hours := g.PlaytimeForever / 60
	mins := g.PlaytimeForever % 60
	return fmt.Sprintf("%s (%dh%02dmin)", g.Name, hours, mins)
}

func SyncTopGames(steamID string, allGames []api.OwnedGame) error {
	sort.Slice(allGames, func(i, j int) bool {
		return allGames[i].PlaytimeForever > allGames[j].PlaytimeForever
	})

	limit := 5
	if len(allGames) < limit {
		limit = len(allGames)
	}
	top := allGames[:limit]

	// Padde à 5 au cas où il y a moins de 5 jeux
	names := make([]string, 5)
	for i := 0; i < limit; i++ {
		names[i] = formatGame(top[i])
	}

	entry := SteamDB{
		SteamID: steamID,
		Game1:   names[0],
		Game2:   names[1],
		Game3:   names[2],
		Game4:   names[3],
		Game5:   names[4],
	}

	// Upsert : met à jour si steam_id existe, sinon insère
	result := db.Where("steam_id = ?", steamID).First(&SteamDB{})
	if result.Error != nil {
		// N'existe pas encore → INSERT
		return db.Create(&entry).Error
	}
	// Existe → UPDATE
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
	if err != nil {
		return nil, err
	}
	return &row, nil
}
