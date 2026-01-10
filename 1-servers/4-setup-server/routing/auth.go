package routing

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/internal/auth"
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
	RespondWithJSON(w, http.StatusOK, MapUserToUserResponse(user))
}
