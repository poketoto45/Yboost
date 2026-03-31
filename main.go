package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"log"

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

	if err := SyncTopGames(steamID, games); err != nil {
		log.Println("Erreur sauvegarde Supabase :", err)
	}

	tpl, err := template.ParseFiles("html/index.html")
	if err != nil {
		http.Error(w, "Fichier HTML introuvable", 500)
		return
	}
	tpl.Execute(w, games)
}
