package routing

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
}

func (config *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
func (config *ApiConfig) HandlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Hits: %d", config.fileserverHits.Load())
}
func (config *ApiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	config.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
