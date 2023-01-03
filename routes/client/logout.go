package client

import (
	"github.com/google/uuid"
	"net/http"
	"vsz-web-backend/database"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil && err != http.ErrNoCookie {
		_, err = uuid.Parse(cookie.Value)
		err := database.DeleteSession(cookie.Value)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		MaxAge: -1,
		Path:   "/",
	})

	http.Redirect(w, r, "/auth/login", http.StatusFound)
}
