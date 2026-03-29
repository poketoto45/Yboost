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
}

func GetOwnedGames(apiKey string, steamID string) ([]OwnedGame, error) {
	// Correction de l'URL (v0001 avec 3 zéros)
	url := fmt.Sprintf("https://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json&include_appinfo=1", apiKey, steamID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur réseau : %v", err)
	}
	defer resp.Body.Close()

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
