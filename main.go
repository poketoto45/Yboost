package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"log"

	"main.go/api" // Ton module s'appelle bien main.go d'après tes messages
)

func main() {
	// 1. Initialisation de la connexion Supabase (défini dans db.go)
	initDB()

	// 2. Configuration du port pour Render
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 3. Gestion des fichiers statiques (CSS)
	// Le dossier "html" contient ton style.css
	fs := http.FileServer(http.Dir("html"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 4. Route principale
	http.HandleFunc("/", homeHandler)

	fmt.Printf("Serveur lancé sur le port %s...\n", port)

	// 5. Lancement du serveur
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Erreur lors du lancement du serveur : ", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// A. Récupération des secrets depuis l'environnement Render
	apiKey := os.Getenv("API_KEY")
	steamID := r.URL.Query().Get("iduser")

	// Si l'utilisateur n'a rien saisi, on utilise l'ID par défaut de Render
	if steamID == "" {
		steamID = os.Getenv("STEAM_ID")
	}

	// Sécurité : on ne continue pas si les variables sont vides
	if apiKey == "" || steamID == "" {
		fmt.Println("Erreur : API_KEY ou STEAM_ID manquant")
		http.Error(w, "Configuration manquante sur le serveur", 500)
		return
	}

	// B. Appel à l'API Steam (défini dans api/steam.go)
	fmt.Printf("Recherche des jeux pour le SteamID : %s\n", steamID)
	games, err := api.GetOwnedGames(apiKey, steamID)

	if err != nil {
		fmt.Printf("Erreur API Steam : %v\n", err)
		http.Error(w, "Impossible de récupérer les jeux : "+err.Error(), 500)
		return
	}

	// C. SAUVEGARDE DANS SUPABASE (CRUD)
	// Cette fonction trie les 5 meilleurs jeux et les enregistre
	err = SyncTopGames(steamID, games)
	if err != nil {
		// On affiche l'erreur en console mais on n'arrête pas le site
		fmt.Println("Erreur lors de la sauvegarde Supabase :", err)
	}

	// D. AFFICHAGE HTML
	tpl, err := template.ParseFiles("html/index.html")
	if err != nil {
		fmt.Println("Erreur template :", err)
		http.Error(w, "Fichier HTML introuvable", 500)
		return
	}

	// On envoie les données au template
	tpl.Execute(w, games)
}
