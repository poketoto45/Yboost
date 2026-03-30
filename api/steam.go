package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Structure pour matcher le JSON de Steam
type SteamResponse struct {
	Response struct {
		Games []OwnedGame `json:"games"`
	} `json:"response"`
}

type OwnedGame struct {
	AppID           int    `json:"appid"`
	Name            string `json:"name"`
	PlaytimeForever int    `json:"playtime_forever"`
	Achievements    []Achievement `json:"achievements"`
}

type Achievement struct {
    APIName  string `json:"apiname"`
    Achieved int    `json:"achieved"` // 1 = débloqué
}

func GetGameAchievements(apiKey string, steamID string, appID int) ([]Achievement, error) {
    url := fmt.Sprintf("https://api.steampowered.com/ISteamUserStats/GetPlayerAchievements/v0001/?appid=%d&key=%s&steamid=%s", appID, apiKey, steamID)

    resp, err := http.Get(url)
    if err != nil || resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("impossible de récupérer les succès")
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

	// SECURITÉ : Si Steam renvoie une erreur (403, 404, 500), on ne parse pas le JSON
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Steam a répondu avec le code %d (Clé API ou ID invalide)", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data SteamResponse
	if err := json.Unmarshal(body, &data); err != nil {
		// C'est ici que ça bloquait. On affiche un bout du body pour débugger au cas où.
		return nil, fmt.Errorf("le JSON est invalide : %v", err)
	}

	return data.Response.Games, nil
}
