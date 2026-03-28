package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"main.go/api" // Vérifie que ton go.mod s'appelle bien "main.go"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// On sert le dossier "html" sous le préfixe "/static/"
	fs := http.FileServer(http.Dir("html"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)

	fmt.Printf("Serveur lancé sur http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("API_KEY")
	// On récupère le steamid depuis le formulaire (?iduser=...) ou depuis l'env
	steamID := r.URL.Query().Get("iduser")
	if steamID == "" {
		steamID = os.Getenv("STEAM_ID") // Valeur par défaut si rien n'est saisi
	}

	// Appel à l'API
	games, err := api.GetOwnedGames(apiKey, steamID)
	if err != nil {
		fmt.Println("Erreur API:", err)
		http.Error(w, "Erreur lors de la récupération des jeux", 500)
		return
	}

	tpl, err := template.ParseFiles("html/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tpl.Execute(w, games)
}
