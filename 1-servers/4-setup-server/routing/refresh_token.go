package routing

import (
	"net/http"

	"example.com/internal/auth"
)

type GetTokenFromRefreshToken struct {
	Token string `json:"token"`
}

// POST /api/refresh
func (config *ApiConfig) HandleCreateTokenByRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	user, err := config.Queries.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	token, err := CreateToken(user.ID, config.JwtSecret)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, GetTokenFromRefreshToken{
		Token: token,
	})
}

// POST /api/revoke
func (config *ApiConfig) HandleRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	err = config.Queries.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	ResponseWithStatus(w, http.StatusNoContent)
}
