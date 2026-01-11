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
type CreateUserResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}
type LoginUserResponse struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
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
		RespondWithError(w, http.StatusInternalServerError, err.Error())
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

// PUT /api/users
func (config *ApiConfig) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	userId, err := CheckValidToken(w, r, config.JwtSecret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	var updateUserRequest createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&updateUserRequest); err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	user, err := config.Queries.GetUserByEmail(r.Context(), updateUserRequest.Email)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if userId != user.ID {
		RespondWithError(w, http.StatusUnauthorized, "cannot update different user")
		return
	}
	password, err := auth.HashPassword(updateUserRequest.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	params := database.UpdateUserParams{
		ID:             userId,
		Email:          updateUserRequest.Email,
		HashedPassword: password,
	}
	updatedUser, err := config.Queries.UpdateUser(r.Context(), params)
	RespondWithJSON(w, http.StatusOK, MapUpdatedUserToUserResponse(updatedUser))
}
func MapUpdatedUserToUserResponse(user database.UpdateUserRow) UserResponse {
	return UserResponse{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
}
func MapUserToUserResponse(user database.User) UserResponse {
	return UserResponse{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
}
func MapUserToLoginUserResponse(user database.User, token string, refreshToken string) LoginUserResponse {
	return LoginUserResponse{
		ID:           user.ID.String(),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}
}
