package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

func GetOwnedGames(apiKey string, steamID string) ([]OwnedGame, error) {
	url := fmt.Sprintf(
		"https://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json&include_appinfo=1",
		apiKey, steamID,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur requête Steam : %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Steam a répondu avec le code %d", resp.StatusCode)
	}

	var data SteamResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("erreur décodage JSON : %w", err)
	}

	if len(data.Response.Games) == 0 {
		return nil, fmt.Errorf("aucun jeu trouvé — profil privé ou SteamID invalide")
	}

	return data.Response.Games, nil
}

func GetGameAchievements(apiKey string, steamID string, appID int) ([]Achievement, error) {
	url := fmt.Sprintf(
		"https://api.steampowered.com/ISteamUserStats/GetPlayerAchievements/v0001/?appid=%d&key=%s&steamid=%s",
		appID, apiKey, steamID,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur requête achievements : %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur ou profil privé")
	}

	var data struct {
		PlayerStats struct {
			Achievements []Achievement `json:"achievements"`
		} `json:"playerstats"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("erreur décodage achievements : %w", err)
	}

	return data.PlayerStats.Achievements, nil
}
