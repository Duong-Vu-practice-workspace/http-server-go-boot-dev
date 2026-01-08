package main

import (
	"fmt"
	"net/http"

	"example.com/routing"
)

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))
	apiConfig := &routing.ApiConfig{}
	appHandler := apiConfig.MiddlewareMetricsInc(http.StripPrefix("/app/", fileServer))
	mux.Handle("/app/", appHandler)
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/metrics", apiConfig.HandlerMetrics)
	mux.HandleFunc("/reset", apiConfig.HandlerReset)
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
