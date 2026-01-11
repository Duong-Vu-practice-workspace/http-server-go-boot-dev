package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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
		log.Fatal(err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))
	apiConfig := &routing.ApiConfig{
		Queries:   dbQueries,
		Db:        db,
		Platform:  os.Getenv("PLATFORM"),
		JwtSecret: os.Getenv("JWT_SECRET"),
	}
	appHandler := apiConfig.MiddlewareMetricsInc(http.StripPrefix("/app/", fileServer))
	mux.Handle("/app/", appHandler)
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /admin/metrics", apiConfig.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiConfig.HandleResetUser)
	mux.HandleFunc("POST /api/users", apiConfig.HandleCreateUser)
	mux.HandleFunc("PUT /api/users", apiConfig.HandleUpdateUser)
	mux.HandleFunc("POST /api/chirps", apiConfig.HandleCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiConfig.HandleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpId}", apiConfig.HandleGetChirpById)
	mux.HandleFunc("DELETE /api/chirps/{chirpId}", apiConfig.HandleDeleteChirpById)
	mux.HandleFunc("POST /api/login", apiConfig.HandleLogin)
	mux.HandleFunc("POST /api/refresh", apiConfig.HandleCreateTokenByRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiConfig.HandleRevokeRefreshToken)
	mux.HandleFunc("POST /api/polka/webhooks", apiConfig.HandlePolkaWebHookMembership)
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
