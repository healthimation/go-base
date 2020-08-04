package goal

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/healthimation/authentication-service/src/token"
	"github.com/healthimation/go-service/alice/middleware"
	"github.com/healthimation/go-service/service"
	"github.com/healthimation/goal-service/src/goal"
	"github.com/healthimation/goal-service/src/internal/data"
	"github.com/husobee/vestigo"
)

func Create(signingKey []byte, dbFactory data.Factory) http.Handler {
	return middleware.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		log := middleware.GetLoggerFromContext(r.Context())
		db, err := dbFactory.Get()
		if err != nil {
			return service.WriteProblem(w, "Service is currently unavailable", goal.ErrorService, http.StatusServiceUnavailable)
		}

		userID := vestigo.Param(r, "userID")
		if userID == "" {
			return service.WriteProblem(w, "user_id must be valid", goal.ErrorInvalidUserID, http.StatusBadRequest)
		}

		// authorization
		jwt := r.Header.Get("Authorization")
		jwt = strings.TrimPrefix(jwt, "Bearer ")
		perm, err := token.VerifyForPerms(jwt, signingKey, log)
		if err != nil {
			return service.WriteProblem(w, "You do not have access to this endpoint.", goal.ErrorNotAuthorized, http.StatusUnauthorized)
		}

		//the user, admin, coach, and cs can call this endpoint
		if !perm.IsUser(userID) && !perm.CanAccessAs([]string{userID}, []string{token.RoleAdmin, token.RoleCoach, token.RoleCS}) {
			return service.WriteProblem(w, "You do not have access to this endpoint.", goal.ErrorNotAuthorized, http.StatusForbidden)
		}

		// the user must be verified
		if !perm.IsVerified() {
			return service.WriteProblem(w, "You do not have access to this endpoint.", goal.ErrorNotAuthorized, http.StatusForbidden)
		}

		req := goal.CreateGoalRequest{}
		dec := json.NewDecoder(r.Body)
		err = dec.Decode(&req)
		if err != nil {
			log.Printf("Error decoding request: %v", err)
			return service.WriteProblem(w, "JSON body must be valid", goal.ErrorInvalidJSON, http.StatusBadRequest)
		}

		if req.StartsAt == nil {
			return service.WriteProblem(w, "starts_at must not be blank", goal.ErrorInvalidGoal, http.StatusBadRequest)
		}

		if req.ExpiresAt == nil {
			return service.WriteProblem(w, "expires_at must not be blank", goal.ErrorInvalidGoal, http.StatusBadRequest)
		}

		if req.ExpiresAt.Before(*req.StartsAt) {
			return service.WriteProblem(w, "expires_at must come after starts_at", goal.ErrorInvalidGoal, http.StatusBadRequest)
		}

		_, startWeek := req.StartsAt.ISOWeek()
		_, expiresWeek := req.ExpiresAt.ISOWeek()

		if req.Category == "weekly" && startWeek != expiresWeek {
			return service.WriteProblem(w, "expires_at and starts_at must be in the same ISO week for weekly goals", goal.ErrorInvalidGoal, http.StatusBadRequest)
		}

		dataErr := db.CreateGoal(r.Context(), userID, req.Type, req.Category, req.Value, req.Name, req.Description, *req.StartsAt, *req.ExpiresAt)
		if dataErr != nil {
			log.Printf("Error getting metrics: %v", dataErr)
			return service.WriteProblem(w, "Error getting metrics.", goal.ErrorService, http.StatusInternalServerError)
		}

		return service.WriteJSONResponse(w, http.StatusOK, `null`)
	})
}
