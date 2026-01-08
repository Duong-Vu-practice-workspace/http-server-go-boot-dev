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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	count := config.fileserverHits.Load()
	s := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
`, count)
	fmt.Fprint(w, s)
}
func (config *ApiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	config.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
