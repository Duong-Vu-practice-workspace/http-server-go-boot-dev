package routing

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"example.com/chirp"
	"example.com/internal/database"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
	Db             *sql.DB
	Queries        *database.Queries
	Platform       string
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

// user type
type createUserRequest struct {
	Email string `json:"email"`
}
type createUserResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

// POST /api/users
func (config *ApiConfig) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		chirp.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Email == "" {
		chirp.RespondWithError(w, http.StatusBadRequest, "email is required")
		return
	}

	user, err := config.Queries.CreateUser(r.Context(), req.Email)
	if err != nil {
		chirp.RespondWithError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	resp := createUserResponse{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	chirp.RespondWithJSON(w, http.StatusCreated, resp)
}

func (config *ApiConfig) HandleResetUser(w http.ResponseWriter, r *http.Request) {
	if config.Platform != "dev" {
		chirp.RespondWithError(w, http.StatusForbidden, "forbidden")
		return
	}
	err := config.Queries.DeleteAllUsers(r.Context())
	if err != nil {
		chirp.RespondWithError(w, http.StatusInternalServerError, "failed to delete users")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
