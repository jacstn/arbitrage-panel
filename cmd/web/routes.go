package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jacstn/arbitrage-panel/internal/handlers"
)

func routes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//mux.Use(WriteToConsole)
	mux.Use(LoadSession)
	mux.Use(NoSurf)
	mux.Get("/", handlers.RunningTrades)
	mux.Get("/all-trades", handlers.AllTrades)
	mux.Get("/about", handlers.About)
	mux.Get("/get-logs/{id}", handlers.GetLogs)
	mux.Get("/market", handlers.Market)
	mux.Get("/list-files", handlers.ListFiles)
	mux.Get("/close-trade/{id}", handlers.CloseTrade)

	mux.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	dataPath := os.Getenv("ARBITRAGE_PY_DIR") + "/data"
	mux.Handle("/data/*", http.StripPrefix("/data", http.FileServer(http.Dir(dataPath))))

	return mux
}
