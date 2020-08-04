package goal

import "time"

type Goal struct {
	ID          int64      `json:"id"`
	UserID      string     `json:"user_id"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	Category    string     `json:"category"`
	Value       int64      `json:"value"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CompletedAt *time.Time `json:"completed_at"`
	SkippedAt   *time.Time `json:"skipped_at"`
}
