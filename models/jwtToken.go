package models

import "time"

type JwtToken struct {
	Token     string     `json:"token"`
	ExpiresAt *time.Time `json:"expiresAt"`
}
