package tenant

type TenantResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Status   string `json:"status"`
	IsActive bool   `json:"is_active"`
}

type CreateTenantRequest struct {
	Name string `json:"name" validate:"required,min=2,max=160"`
	Slug string `json:"slug" validate:"required,min=2,max=120"`
}

type UpdateTenantRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=160"`
	Status   string `json:"status" validate:"required"`
	IsActive bool   `json:"is_active"`
}
