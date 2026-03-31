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

	sql := `
		INSERT INTO "steamDB" (steam_id, game1, game2, game3, game4, game5)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (steam_id) DO UPDATE SET
			game1 = EXCLUDED.game1,
			game2 = EXCLUDED.game2,
			game3 = EXCLUDED.game3,
			game4 = EXCLUDED.game4,
			game5 = EXCLUDED.game5
	`
	return db.Exec(sql, steamID, names[0], names[1], names[2], names[3], names[4]).Error
}

func GetTopGamesFromDB(steamID string) (*SteamDB, error) {
	var row SteamDB
	err := db.Raw(`SELECT * FROM "steamDB" WHERE steam_id = ? LIMIT 1`, steamID).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.SteamID == "" {
		return nil, fmt.Errorf("aucune donnée pour ce joueur")
	}
	return &row, nil
}

func DeleteTopGames(steamID string) error {
	result := db.Exec(`DELETE FROM "steamDB" WHERE steam_id = ?`, steamID)
	if result.RowsAffected == 0 {
		return fmt.Errorf("aucune donnée trouvée pour ce joueur")
	}
	return result.Error
}
