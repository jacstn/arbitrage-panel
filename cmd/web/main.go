package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/jacstn/arbitrage-panel/config"
	"github.com/jacstn/arbitrage-panel/internal/database"
	"github.com/jacstn/arbitrage-panel/internal/handlers"
)

const portNumber = ":3333"

var app = config.AppConfig{}

func main() {
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Secure = app.Production

	db := database.Connect()
	app.DB = db

	app.Session = session
	app.PythonScriptsDir = os.Getenv("ARBITRAGE_PY_DIR")
	app.PythonExecuteCmd = os.Getenv("ARBITRAGE_PY")
	if app.PythonExecuteCmd == "" || app.PythonScriptsDir == "" {
		log.Fatal("Cannot start server, python envirnoment error")
	}

	if os.Getenv("RUN_MODE") == "PRODUCTION" {
		app.Production = true
	} else {
		app.Production = false
	}

	handlers.NewHandlers(&app)
	fmt.Println("Starting application", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}
	defer db.Close()
	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal("Cannot start server", err)
	}
}
