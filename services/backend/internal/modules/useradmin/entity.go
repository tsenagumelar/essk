package useradmin

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	TenantID *uuid.UUID
	Email    string
	Name     string
	Status   string
	IsActive bool
}
