package rbac

import "github.com/google/uuid"

type Permission struct {
	ID          uuid.UUID
	Code        string
	Name        string
	Description *string
	IsActive    bool
}

type Role struct {
	ID          uuid.UUID
	TenantID    *uuid.UUID
	Name        string
	Code        string
	Description *string
	IsSystem    bool
	IsActive    bool
}
