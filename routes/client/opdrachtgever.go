package client

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"vsz-web-backend"
	"vsz-web-backend/database"
)

func GetOpdrachtgever(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*vsz_web_backend.Opdrachtgever)

	opdrachtgever, err := database.GetOpdrachtgeverByID(user.Bedrijfscode)
	if err != nil {
		log.Printf("failed to fetch opdrachtgever: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(opdrachtgever)
	if err != nil {
		log.Printf("failed to marshal opdrachtgever: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func GetOpdrachtgeverCount(w http.ResponseWriter, r *http.Request) {
	opdrachtgevers, err := database.GetOpdrachtgeverCount()
	if err != nil {
		log.Printf("failed to fetch opdrachtgevers: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(strconv.Itoa(opdrachtgevers)))
}
