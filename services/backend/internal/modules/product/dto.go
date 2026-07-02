package product

type ProductResponse struct {
	ID         string `json:"id"`
	TenantID   string `json:"tenant_id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	PriceCents int    `json:"price_cents"`
	Status     string `json:"status"`
	IsActive   bool   `json:"is_active"`
}

type CreateProductRequest struct {
	Code       string `json:"code" validate:"required,min=2,max=80"`
	Name       string `json:"name" validate:"required,min=2,max=160"`
	Category   string `json:"category" validate:"required,min=2,max=120"`
	PriceCents int    `json:"price_cents" validate:"gte=0"`
}

type UpdateProductRequest struct {
	Name       string `json:"name" validate:"required,min=2,max=160"`
	Category   string `json:"category" validate:"required,min=2,max=120"`
	PriceCents int    `json:"price_cents" validate:"gte=0"`
	Status     string `json:"status" validate:"required,oneof=active inactive draft"`
	IsActive   bool   `json:"is_active"`
}
