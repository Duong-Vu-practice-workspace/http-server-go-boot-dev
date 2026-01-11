package routing

import (
	"encoding/json"
	"net/http"
	"time"

	"example.com/internal/auth"
	"example.com/internal/database"
	"github.com/google/uuid"
)

type createChirpRequest struct {
	Body   string    `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}
type chirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

const CHIRP_ID = "chirpId"

// POST /api/chirps.sql
func (config *ApiConfig) HandleCreateChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	userId, err := auth.ValidateJWT(token, config.JwtSecret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	var request createChirpRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	body := request.Body
	bodyUserId := request.UserId
	if bodyUserId != userId {
		RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if body == "" {
		RespondWithError(w, http.StatusBadRequest, "Body is required")
		return
	}
	if userId == uuid.Nil {
		RespondWithError(w, http.StatusBadRequest, "User ID is required")
		return
	}
	validateChirp := ValidateChirp{Body: body}
	result, err := ValidateChirpHandler(w, r, validateChirp)
	if err != nil {
		return
	}
	params := database.CreateChirpParams{
		Body:   result,
		UserID: userId,
	}
	createdChirp, err := config.Queries.CreateChirp(r.Context(), params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "failed to create chirp")
		return
	}

	RespondWithJSON(w, http.StatusCreated, mapCreateChirpToResponse(createdChirp))

}
func (config *ApiConfig) HandleGetChirps(w http.ResponseWriter, r *http.Request) {
	res, err := config.Queries.GetAllChirps(r.Context())
	result := make([]chirpResponse, len(res))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "failed to get all chirps")
		return
	}
	for i := range res {
		result[i] = mapCreateChirpToResponse(res[i])
	}
	RespondWithJSON(w, http.StatusOK, result)
}
func (config *ApiConfig) HandleGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpIdString := r.PathValue(CHIRP_ID)
	chirpId, err := uuid.Parse(chirpIdString)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid chirp id")
		return
	}

	chirp, err := config.Queries.GetChirpById(r.Context(), chirpId)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "failed to get all chirps")
		return
	}
	RespondWithJSON(w, http.StatusOK, mapCreateChirpToResponse(chirp))
}
func mapCreateChirpToResponse(createdChirp database.Chirp) chirpResponse {
	return chirpResponse{
		ID:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.CreatedAt,
		Body:      createdChirp.Body,
		UserID:    createdChirp.UserID,
	}
}
