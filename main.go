package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	api "main.go/api"
)

type game struct {
	ID                            int    `json:"id"`
	Name                          string `json:"name"`
	Playtime_2weeks               int    `json:"playtime_2weeks"`
	Playtime_forever              int    `json:"playtime_forever"`
	Img_icon_url                  string `json:"img_icon_url"`
	Has_community_visible_stats   bool   `json:"has_community_visible_stats"`
	Playtime_windows_forever      int    `json:"playtime_windows_forever"`
	Playtime_mac_forever          int    `json:"playtime_mac_forever"`
	Playtime_linux_forever        int    `json:"playtime_linux_forever"`
	Playtime_deck_forever         int    `json:"playtime_deck_forever"`
	Rtime_last_played             int    `json:"rtime_last_played"`
	Playtime_disconnected_forever int    `json:"playtime_disconnected_forever"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// CORRECTION : Ton dossier s'appelle "html", pas "front/static"
	// On sert tout le dossier "html" pour le CSS
	fs := http.FileServer(http.Dir("html"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)

	fmt.Printf("Serveur lancé sur http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Erreur : %v\n", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	body, err := api.ChooseArtisteAll()
	if err != nil {
		http.Error(w, "Erreur API", 500)
		return
	}

	// Debug : affiche le début de ce que reçoit Go dans les logs Render
	if len(body) > 0 {
		fmt.Printf("Réponse reçue (50 premiers caractères): %s\n", string(body[:50]))
	}

	var users []game
	marshall1(body, &users)

	// CORRECTION : Le chemin doit être exact
	tpl, err := template.ParseFiles("html/index.html")
	if err != nil {
		http.Error(w, "Template introuvable : "+err.Error(), 500)
		return
	}

	tpl.Execute(w, users)
}

func marshall1(jsonFromApi []byte, truc *[]game) {

	err := json.Unmarshal(jsonFromApi, truc)
	if err != nil {
		fmt.Println("Error unmarshalling json:", err)
	}
}

func marshall2(jsonFromApi []byte, truc *game) {

	err := json.Unmarshal(jsonFromApi, truc)
	if err != nil {
		fmt.Println("Error unmarshalling json:", err)
	}
}
