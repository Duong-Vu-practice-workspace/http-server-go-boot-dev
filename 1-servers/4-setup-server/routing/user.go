package routing

import (
	"encoding/json"
	"net/http"
	"time"

	"example.com/internal/auth"
	"example.com/internal/database"
)

// user type
type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserResponse struct {
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
	password, err := auth.HashPassword(req.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	request := database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: password,
	}
	user, err := config.Queries.CreateUser(r.Context(), request)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	RespondWithJSON(w, http.StatusCreated, MapUserToUserResponse(user))
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
func MapUserToUserResponse(user database.User) UserResponse {
	return UserResponse{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
}
