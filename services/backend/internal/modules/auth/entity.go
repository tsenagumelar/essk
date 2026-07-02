package auth

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	TenantID     *uuid.UUID
	Email        string
	Name         string
	PasswordHash string
	Status       string
	IsActive     bool
	IsDeleted    bool
}

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
}
