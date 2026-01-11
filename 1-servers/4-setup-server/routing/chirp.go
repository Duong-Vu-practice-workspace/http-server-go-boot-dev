package routing

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

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

const ChirpId = "chirpId"
const AuthorId = "author_id"
const Sort = "sort"

// POST /api/chirps.sql
func (config *ApiConfig) HandleCreateChirp(w http.ResponseWriter, r *http.Request) {
	userId, err := CheckValidToken(w, r, config.JwtSecret)
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
	authorIdString := r.URL.Query().Get(AuthorId)
	sortString := r.URL.Query().Get(Sort)
	var err error
	var res []database.Chirp

	if authorIdString == "" {
		res, err = config.Queries.GetAllChirps(r.Context())
	} else {
		authorId, err := uuid.Parse(authorIdString)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "invalid author id")
			return
		}
		res, err = config.Queries.GetAllChirpsByAuthorSortCreatedAtAscending(r.Context(), authorId)
	}
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "failed to get all chirps")
		return
	}
	result := make([]chirpResponse, len(res))

	for i := range res {
		result[i] = mapCreateChirpToResponse(res[i])
	}
	sort.Slice(result, func(i, j int) bool {
		if sortString == "desc" {
			return result[i].CreatedAt.After(result[j].CreatedAt)
		}
		//default asc
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	RespondWithJSON(w, http.StatusOK, result)
}
func (config *ApiConfig) HandleGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpIdString := r.PathValue(ChirpId)
	chirpId, err := uuid.Parse(chirpIdString)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid chirp id")
		return
	}

	chirp, err := config.Queries.GetChirpById(r.Context(), chirpId)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "failed to get chirp by Id")
		return
	}
	RespondWithJSON(w, http.StatusOK, mapCreateChirpToResponse(chirp))
}
func (config *ApiConfig) HandleDeleteChirpById(w http.ResponseWriter, r *http.Request) {
	userId, err := CheckValidToken(w, r, config.JwtSecret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	chirpIdString := r.PathValue(ChirpId)
	chirpId, err := uuid.Parse(chirpIdString)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid chirp id")
		return
	}

	chirp, err := config.Queries.GetChirpById(r.Context(), chirpId)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "failed to get chirp by Id")
		return
	}
	if chirp.UserID != userId {
		RespondWithError(w, http.StatusForbidden, "you cannot delete other user's chirp")
		return
	}
	_, err = config.Queries.DeleteChirpById(r.Context(), chirpId)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ResponseWithStatus(w, http.StatusNoContent)
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
