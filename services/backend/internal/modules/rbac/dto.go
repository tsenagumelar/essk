package rbac

type PermissionResponse struct {
	ID          string  `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	IsActive    bool    `json:"is_active"`
}

type RoleResponse struct {
	ID          string  `json:"id"`
	TenantID    *string `json:"tenant_id,omitempty"`
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Description *string `json:"description,omitempty"`
	IsSystem    bool    `json:"is_system"`
	IsActive    bool    `json:"is_active"`
}

type CreateRoleRequest struct {
	TenantID    *string `json:"tenant_id"`
	Name        string  `json:"name" validate:"required,min=2,max=120"`
	Code        string  `json:"code" validate:"required,min=2,max=120"`
	Description *string `json:"description"`
	IsSystem    bool    `json:"is_system"`
}

type UpdateRoleRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=120"`
	Description *string `json:"description"`
	IsActive    bool    `json:"is_active"`
}

type AssignPermissionRequest struct {
	PermissionID string `json:"permission_id" validate:"required,uuid"`
}

type AssignRoleRequest struct {
	RoleID string `json:"role_id" validate:"required,uuid"`
}
