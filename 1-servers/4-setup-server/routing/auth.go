package routing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/internal/auth"
	"example.com/internal/database"
	"github.com/google/uuid"
)

type loginUser struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (config *ApiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	user, err := config.Queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("cannot find user with email = %s", req.Email))
		return
	}
	isMatch, err := auth.CheckPasswordHash(req.Password, user.HashedPassword)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error validating password")
		return
	}
	if !isMatch {
		RespondWithError(w, http.StatusUnauthorized, "password not match")
		return
	}

	token, err := CreateToken(user.ID, config.JwtSecret)
	var refreshToken string
	refreshToken, err = config.Queries.GetRefreshTokenFromUser(r.Context(), user.ID)
	if err != nil {
		refreshToken, err = auth.MakeRefreshToken()
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "error creating refresh token")
			return
		}
		params := database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(time.Minute * 60 * 24 * 60),
		}
		_, err = config.Queries.CreateRefreshToken(r.Context(), params)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	RespondWithJSON(w, http.StatusOK, MapUserToLoginUserResponse(user, token, refreshToken))
}
func CreateToken(userId uuid.UUID, tokenSecret string) (string, error) {
	expiresInSeconds := 3600 //1hour
	return auth.MakeJWT(userId, tokenSecret, time.Duration(expiresInSeconds)*time.Second)
}
func CheckValidToken(w http.ResponseWriter, r *http.Request, tokenSecret string) (uuid.UUID, error) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		//RespondWithError(w, http.StatusUnauthorized, err.Error())
		return uuid.Nil, err
	}
	userId, err := auth.ValidateJWT(token, tokenSecret)
	if err != nil {
		//RespondWithError(w, http.StatusUnauthorized, err.Error())
		return uuid.Nil, err
	}
	return userId, nil
}
