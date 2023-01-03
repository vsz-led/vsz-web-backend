package auth

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
	"vsz-web-backend"
	"vsz-web-backend/database"
)

func Login(w http.ResponseWriter, r *http.Request) {
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	if password == "" || email == "" {
		http.Error(w, "The supplied password/email cannot be empty.", http.StatusBadRequest)
		return
	}

	user, err := database.GetOpdrachtgeverByEmail(email)
	if err != nil {
		log.Printf("Failed to get user by email: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "The email and password combination is incorrect.", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Wachtwoord, []byte(password))
	if err != nil {
		http.Error(w, "The email and password combination is incorrect.", http.StatusBadRequest)
		return
	}

	session := &vsz_web_backend.Session{
		User: user.Bedrijfscode,
	}

	sessionuuid, err := database.StoreSession(session)
	if err != nil {
		log.Printf("Failed to store session: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionuuid.String(),
		Expires: time.Now().Add(database.SessionTime),
		Path:    "/",
	})

	http.Redirect(w, r, "/client", http.StatusFound)
}
