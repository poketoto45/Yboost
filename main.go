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
	fmt.Println("Hello, World!")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	})

	http.ListenAndServe(":"+port, nil)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("front/static"))))

	// Lier la page d'accueil à ton handler
	http.HandleFunc("/home", homeHandler)

	http.ListenAndServe(":"+port, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	body, err := api.ChooseArtisteAll()
	var users []game
	marshall1(body, &users)

	tpl, err := template.ParseFiles("front/page/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
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
