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

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/sync", syncHandler)
	http.HandleFunc("/top", topGamesHandler)

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
		http.Error(w, "Configuration manquante", 500)
		return
	}

	games, err := api.GetOwnedGames(apiKey, steamID)
	if err != nil {
		http.Error(w, "Impossible de récupérer les jeux : "+err.Error(), 500)
		return
	}

	tpl, err := template.ParseFiles("html/index.html")
	if err != nil {
		http.Error(w, "Fichier HTML introuvable", 500)
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
	if steamID == "" {
		steamID = os.Getenv("STEAM_ID")
	}

	if apiKey == "" || steamID == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "steamID ou API_KEY manquant"})
		return
	}

	games, err := api.GetOwnedGames(apiKey, steamID)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := SyncTopGames(steamID, games); err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "Top 5 sauvegardé dans Supabase !",
	})
}

func topGamesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	steamID := r.URL.Query().Get("iduser")
	if steamID == "" {
		steamID = os.Getenv("STEAM_ID")
	}

	row, err := GetTopGamesFromDB(steamID)
	if err != nil {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]string{"error": "Aucune donnée pour ce joueur"})
		return
	}

	json.NewEncoder(w).Encode(row)
}
