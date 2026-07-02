package product

import "github.com/google/uuid"

type Product struct {
	ID         uuid.UUID
	TenantID   uuid.UUID
	Code       string
	Name       string
	Category   string
	PriceCents int
	Status     string
	IsActive   bool
}
