package client

import (
	"encoding/json"
	"log"
	"net/http"
	"vsz-web-backend/database"
)

func GetKruisingen(w http.ResponseWriter, r *http.Request) {
	kruisingen, err := database.GetKruisingen()
	if err != nil {
		log.Printf("failed to fetch kruisingen: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(kruisingen)
	if err != nil {
		log.Printf("failed to marshal kruisingen: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
