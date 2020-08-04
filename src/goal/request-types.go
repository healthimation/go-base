package goal

import "time"

type CreateGoalRequest struct {
	Type        string     `json:"type"`
	Category    string     `json:"category"`
	Value       int64      `json:"value"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	StartsAt    *time.Time `json:"starts_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
}
