package auth

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"log"
	"net/http"
	"vsz-web-backend"
	"vsz-web-backend/config"
	"vsz-web-backend/database"
)

func CheckLoggedIn() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			user, err := getCookie(w, r)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				if config.Global.Debug {
					log.Printf("Encountered error when getting user session: %s", err)
				}
				return
			}
			if user == nil {
				http.Redirect(w, r, "/auth/login", http.StatusFound)
				return
			}

			requestcontext := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(requestcontext))
		}
		return http.HandlerFunc(fn)
	}
}

func CheckNotLoggedIn() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			user, err := getCookie(w, r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			if user == nil {
				next.ServeHTTP(w, r)
				return
			}

			http.Redirect(w, r, "/client", http.StatusFound)
		}
		return http.HandlerFunc(fn)
	}
}

func getCookie(w http.ResponseWriter, r *http.Request) (*vsz_web_backend.Opdrachtgever, error) {
	cookie, err := r.Cookie("session_token")
	if err == http.ErrNoCookie {
		return nil, nil
	}

	_, err = uuid.Parse(cookie.Value)
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			MaxAge: -1,
			Path:   "/",
		})
		return nil, nil
	}

	session, err := database.GetSession(cookie.Value)
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			MaxAge: -1,
			Path:   "/",
		})
		return nil, err
	}

	if session == nil {
		return nil, nil
	}

	user, err := database.GetOpdrachtgeverByID(session.User)
	if err != nil {
		return nil, err
	}
	if user == nil {
		err = database.DeleteSession(cookie.Value)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("unknown user")
	}

	return user, nil
}
