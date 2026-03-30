package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 1. Les structures en haut
type OwnedGame struct {
	AppID           int           `json:"appid"`
	Name            string        `json:"name"`
	PlaytimeForever int           `json:"playtime_forever"`
	Achievements    []Achievement `json:"achievements"`
}

type Achievement struct {
	APIName  string `json:"apiname"`
	Achieved int    `json:"achieved"`
}

type SteamResponse struct {
	Response struct {
		Games []OwnedGame `json:"games"`
	} `json:"response"`
}

// 2. Ta fonction GetOwnedGames
func GetOwnedGames(apiKey string, steamID string) ([]OwnedGame, error) {
	url := fmt.Sprintf("https://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json&include_appinfo=1", apiKey, steamID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data SteamResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Response.Games, nil
}

// 3. Ta nouvelle fonction pour les succès
func GetGameAchievements(apiKey string, steamID string, appID int) ([]Achievement, error) {
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUserStats/GetPlayerAchievements/v0001/?appid=%d&key=%s&steamid=%s", appID, apiKey, steamID)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur ou profil privé")
	}
	defer resp.Body.Close()

	var data struct {
		PlayerStats struct {
			Achievements []Achievement `json:"achievements"`
		} `json:"playerstats"`
	}

	json.NewDecoder(resp.Body).Decode(&data)
	return data.PlayerStats.Achievements, nil
}
