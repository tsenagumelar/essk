package useradmin

type UserResponse struct {
	ID       string   `json:"id"`
	TenantID *string  `json:"tenant_id,omitempty"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Status   string   `json:"status"`
	IsActive bool     `json:"is_active"`
	RoleIDs  []string `json:"role_ids"`
}

type CreateUserRequest struct {
	TenantID string   `json:"tenant_id" validate:"required,uuid"`
	Email    string   `json:"email" validate:"required,email,max=255"`
	Name     string   `json:"name" validate:"required,min=2,max=160"`
	Password string   `json:"password" validate:"required,min=8,max=128"`
	RoleIDs  []string `json:"role_ids"`
}

type UpdateUserRequest struct {
	Name     string   `json:"name" validate:"required,min=2,max=160"`
	Status   string   `json:"status" validate:"required,oneof=active inactive invited suspended"`
	IsActive bool     `json:"is_active"`
	RoleIDs  []string `json:"role_ids"`
}
