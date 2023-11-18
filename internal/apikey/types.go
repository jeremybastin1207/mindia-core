package apikey

import "time"

type ApiKeyMap = map[string]ApiKey

type ApiKey struct {
	Name      string    `json:"name,omitempty"`
	Key       string    `json:"key,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
