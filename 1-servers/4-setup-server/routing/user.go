package routing

import (
	"encoding/json"
	"net/http"
	"time"
)

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
		RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Email == "" {
		RespondWithError(w, http.StatusBadRequest, "email is required")
		return
	}

	user, err := config.Queries.CreateUser(r.Context(), req.Email)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	resp := createUserResponse{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	RespondWithJSON(w, http.StatusCreated, resp)
}

func (config *ApiConfig) HandleResetUser(w http.ResponseWriter, r *http.Request) {
	if config.Platform != "dev" {
		RespondWithError(w, http.StatusForbidden, "forbidden")
		return
	}
	err := config.Queries.DeleteAllUsers(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "failed to delete users")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
