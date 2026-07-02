package tenant

import "github.com/google/uuid"

type Tenant struct {
	ID       uuid.UUID
	Name     string
	Slug     string
	Status   string
	IsActive bool
}
