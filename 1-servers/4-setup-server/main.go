package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"example.com/chirp"
	"example.com/internal/database"
	"example.com/routing"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		// handle error appropriately (log and exit)
		panic(err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))
	apiConfig := &routing.ApiConfig{}
	appHandler := apiConfig.MiddlewareMetricsInc(http.StripPrefix("/app/", fileServer))
	mux.Handle("/app/", appHandler)
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /admin/metrics", apiConfig.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiConfig.HandlerReset)
	mux.HandleFunc("POST /api/validate_chirp", chirp.ValidateChirpHandler)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Errorf("server error :%v", err)
	}
}
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
