package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// SteamResponse correspond à la hiérarchie réelle de l'API Steam
type SteamResponse struct {
	Response struct {
		Games []OwnedGame `json:"games"`
	} `json:"response"`
}

type OwnedGame struct {
	AppID           int    `json:"appid"`
	Name            string `json:"name"`
	PlaytimeForever int    `json:"playtime_forever"`
	ImgIconURL      string `json:"img_icon_url"`
	Playtime2Weeks  int    `json:"playtime_2weeks"`

}

func GetOwnedGames(apiKey string, steamID string) ([]OwnedGame, error) {
	url := fmt.Sprintf("https://api.steampowered.com/IPlayerService/GetOwnedGames/v001/?key=%s&steamid=%s&format=json&include_appinfo=1", apiKey, steamID)
	fmt.Println("Appel de l'API pour le SteamID:", steamID) // Ne pas afficher la clé pour la sécurité

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var data SteamResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data.Response.Games, nil
}
