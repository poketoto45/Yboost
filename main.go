package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"main.go/api" // Vérifie que ton module s'appelle bien "main.go"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Servir le CSS
	fs := http.FileServer(http.Dir("html"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)

	fmt.Printf("Serveur prêt sur le port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Récupération des variables
	apiKey := os.Getenv("API_KEY")
	steamID := r.URL.Query().Get("iduser")

	// Fallback sur le STEAM_ID par défaut de Render si l'input est vide
	if steamID == "" {
		steamID = os.Getenv("STEAM_ID")
	}

	// 2. Vérification avant l'appel
	if apiKey == "" || steamID == "" {
		fmt.Println("Erreur : API_KEY ou STEAM_ID manquant dans Render")
		http.Error(w, "Configuration manquante (Clé ou ID)", 500)
		return
	}

	// 3. Appel API
	fmt.Printf("Tentative d'appel pour l'ID: %s\n", steamID)
	games, err := api.GetOwnedGames(apiKey, steamID)

	if err != nil {
		fmt.Printf("ERREUR : %v\n", err)
		http.Error(w, "Erreur lors de la récupération : "+err.Error(), 500)
		return
	}

	// 4. Affichage
	tpl, err := template.ParseFiles("html/index.html")
	if err != nil {
		http.Error(w, "Fichier HTML introuvable", 500)
		return
	}

	tpl.Execute(w, games)
}
