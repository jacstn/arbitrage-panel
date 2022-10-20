package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/jacstn/arbitrage-panel/config"
	"github.com/jacstn/arbitrage-panel/internal/database"
	"github.com/jacstn/arbitrage-panel/internal/handlers"
)

const portNumber = ":3000"

var app = config.AppConfig{
	Production: false,
}

func main() {
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Secure = app.Production

	db := database.Connect()
	app.DB = db

	app.Session = session
	handlers.NewHandlers(&app)
	fmt.Println("Starting application", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}
	defer db.Close()
	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal("Cannot start server")
	}
}
