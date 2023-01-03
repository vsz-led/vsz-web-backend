package client

import (
	"encoding/json"
	"log"
	"net/http"
	"vsz-web-backend/database"
)

func GetAutos(w http.ResponseWriter, r *http.Request) {
	autos, err := database.GetAutos()
	if err != nil {
		log.Printf("failed to fetch autos: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(autos)
	if err != nil {
		log.Printf("failed to marshal autos: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
