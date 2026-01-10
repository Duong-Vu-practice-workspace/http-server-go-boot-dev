package routing

import (
	"encoding/json"
	"net/http"
	"time"

	"example.com/internal/database"
	"github.com/google/uuid"
)

type createChirpRequest struct {
	Body   string    `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}
type createChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

// POST /api/chirps.sql
func (config *ApiConfig) HandleCreateChirp(w http.ResponseWriter, r *http.Request) {
	var request createChirpRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	body := request.Body
	userId := request.UserId
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

	responseChirp := createChirpResponse{
		ID:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.CreatedAt,
		Body:      createdChirp.Body,
		UserID:    createdChirp.UserID,
	}
	RespondWithJSON(w, http.StatusCreated, responseChirp)

}
