package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func ChooseArtisteAll() ([]byte, error) {

	url := os.Getenv("STEAM_API_URL")
	response, err := http.Get(url)

	if err != nil {

		fmt.Println("Erreur lors de la requête HTTP :", err)

		return nil, err

	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {

		fmt.Println("Erreur lors de la lecture de la réponse :", err)

		return nil, err

	}

	return body, nil
}
