package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"log"
	"net/http"
	"time"
	"vsz-web-backend/config"
	"vsz-web-backend/database"
	"vsz-web-backend/routes/auth"
	"vsz-web-backend/routes/client"
)

func main() {
	// load configuration file (config.yml)
	err := config.Load()
	if err != nil {
		log.Printf("failed to load config: %s", err)
		return
	}

	// load database
	err = database.Initialize()
	if err != nil {
		log.Printf("failed to initialize database: %s", err)
		return
	}

	// initialize http router
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	if config.Global.Debug {
		log.Println("Starting up in debug mode")
		r.Use(middleware.Logger)
	}

	// create route for authentication
	r.Route("/auth", func(r chi.Router) {
		r.Use(httprate.LimitByIP(10, time.Minute))
		r.Use(auth.CheckNotLoggedIn())
		r.Post("/login", auth.Login)
	})

	// create route for client
	r.Route("/client", func(r chi.Router) {
		//r.Use(auth.CheckLoggedIn())
		r.Get("/opdrachtgever", client.GetOpdrachtgever)
		r.Get("/opdrachtgevercount", client.GetOpdrachtgeverCount)
		r.Get("/kruisingen", client.GetKruisingen)
		r.Get("/autos", client.GetAutos)
		r.Get("/autosweek", client.GetAutosWeek)
		r.Get("/autosmaand", client.GetAutosMaand)
		r.Get("/autoskruising", client.GetAutosKruising)
		r.Get("/autosopdrachtgever", client.GetAutosOpdrachtgever)
		r.Post("/logout", client.Logout)
	})

	// start http listener for router
	l := config.Global.Listen
	log.Printf("Starting web listener on %s", l)
	err = http.ListenAndServe(l, r)
	if err != nil {
		log.Printf("Failed to start web listener: %s", err)
		return
	}
}
