package routing

import (
	"encoding/json"
	"net/http"

	"example.com/internal/database"
	"github.com/google/uuid"
)

type UserUpgradeMembershipData struct {
	UserId uuid.UUID `json:"user_id"`
}
type UserUpgradeMembershipWebhook struct {
	Event string                    `json:"event"`
	Data  UserUpgradeMembershipData `json:"data"`
}

const UpgradeMembershipEvent = "user.upgraded"

func (config *ApiConfig) HandlePolkaWebHookMembership(w http.ResponseWriter, r *http.Request) {
	var request UserUpgradeMembershipWebhook
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ResponseWithStatus(w, http.StatusBadRequest)
		return
	}
	if request.Event != UpgradeMembershipEvent {
		ResponseWithStatus(w, http.StatusNoContent)
		return
	}
	params := database.UpdateUserMemberShipParams{
		ID:          request.Data.UserId,
		IsChirpyRed: true,
	}
	_, err := config.Queries.UpdateUserMemberShip(r.Context(), params)
	if err != nil {
		ResponseWithStatus(w, http.StatusNotFound)
		return
	}
	ResponseWithStatus(w, http.StatusNoContent)
}
