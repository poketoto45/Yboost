package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"main.go/api"
)

func main() {
	initDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fs := http.FileServer(http.Dir("html"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/sync", syncHandler)
	http.HandleFunc("/top", topGamesHandler)
	http.HandleFunc("/delete", deleteHandler)

	fmt.Printf("Serveur lancé sur le port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Erreur serveur : ", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("API_KEY")
	steamID := r.URL.Query().Get("iduser")
	if steamID == "" {
		steamID = os.Getenv("STEAM_ID")
	}

	if apiKey == "" || steamID == "" {
		http.Error(w, "Configuration API_KEY ou STEAM_ID manquante", 500)
		return
	}

	games, err := api.GetOwnedGames(apiKey, steamID)
	if err != nil {
		renderTemplate(w, "html/index.html", nil, steamID)
		return
	}

	renderTemplate(w, "html/index.html", games, steamID)
}

func renderTemplate(w http.ResponseWriter, path string, games []api.OwnedGame, steamID string) {
	tpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, "Erreur template", 500)
		return
	}
	data := struct {
		Games   []api.OwnedGame
		SteamID string
	}{
		Games:   games,
		SteamID: steamID,
	}
	tpl.Execute(w, data)
}

func syncHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	apiKey := os.Getenv("API_KEY")
	steamID := r.URL.Query().Get("iduser")

	games, err := api.GetOwnedGames(apiKey, steamID)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := SyncTopGames(steamID, games); err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erreur DB: " + err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Top 5 synchronisé !"})
}

func topGamesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	steamID := r.URL.Query().Get("iduser")

	row, err := GetTopGamesFromDB(steamID)
	if err != nil {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]string{"error": "Données introuvables"})
		return
	}
	json.NewEncoder(w).Encode(row)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	steamID := r.URL.Query().Get("iduser")

	if err := DeleteTopGames(steamID); err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Données supprimées avec succès"})
}
