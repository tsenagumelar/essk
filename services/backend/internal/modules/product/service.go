package product

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	apperrors "github.com/tsenagumelar/essk/services/backend/internal/errors"
	"github.com/tsenagumelar/essk/services/backend/internal/modules/audit"
)

type Service struct {
	repo  Repository
	audit *audit.Service
	now   func() time.Time
}

func NewService(repo Repository) Service {
	return Service{repo: repo, now: time.Now}
}

func (s Service) WithAudit(auditService audit.Service) Service {
	s.audit = &auditService
	return s
}

func (s Service) List(ctx context.Context, tenantID uuid.UUID) ([]ProductResponse, error) {
	products, err := s.repo.List(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	result := make([]ProductResponse, 0, len(products))
	for _, product := range products {
		result = append(result, toResponse(product))
	}
	return result, nil
}

func (s Service) Get(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (ProductResponse, error) {
	product, err := s.repo.Get(ctx, tenantID, id)
	if err != nil {
		return ProductResponse{}, mapNotFound(err)
	}
	return toResponse(product), nil
}

func (s Service) Create(ctx context.Context, tenantID uuid.UUID, req CreateProductRequest, actorID uuid.UUID) (ProductResponse, error) {
	product := Product{
		ID:         uuid.New(),
		TenantID:   tenantID,
		Code:       req.Code,
		Name:       req.Name,
		Category:   req.Category,
		PriceCents: req.PriceCents,
		Status:     "active",
		IsActive:   true,
	}
	if err := s.repo.Create(ctx, product, actorID, s.now().UTC()); err != nil {
		return ProductResponse{}, err
	}
	_ = s.writeAudit(ctx, actorID, "product.create", product.ID.String(), map[string]any{"code": product.Code})
	return toResponse(product), nil
}

func (s Service) Update(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, req UpdateProductRequest, actorID uuid.UUID) (ProductResponse, error) {
	product, err := s.repo.Get(ctx, tenantID, id)
	if err != nil {
		return ProductResponse{}, mapNotFound(err)
	}
	product.Name = req.Name
	product.Category = req.Category
	product.PriceCents = req.PriceCents
	product.Status = req.Status
	product.IsActive = req.IsActive
	if err := s.repo.Update(ctx, product, actorID, s.now().UTC()); err != nil {
		return ProductResponse{}, mapNotFound(err)
	}
	_ = s.writeAudit(ctx, actorID, "product.update", product.ID.String(), map[string]any{"status": product.Status})
	return toResponse(product), nil
}

func (s Service) Delete(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID uuid.UUID) error {
	if err := mapNotFound(s.repo.Delete(ctx, tenantID, id, actorID, s.now().UTC())); err != nil {
		return err
	}
	_ = s.writeAudit(ctx, actorID, "product.delete", id.String(), nil)
	return nil
}

func (s Service) writeAudit(ctx context.Context, actorID uuid.UUID, action string, resourceID string, metadata map[string]any) error {
	if s.audit == nil {
		return nil
	}
	return s.audit.Write(ctx, audit.Event{
		ActorUserID:  &actorID,
		Action:       action,
		ResourceType: "product",
		ResourceID:   &resourceID,
		Metadata:     metadata,
	})
}

func mapNotFound(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrNotFound) {
		return apperrors.New("NOT_FOUND", fiber.StatusNotFound, "Resource not found")
	}
	return err
}

func toResponse(product Product) ProductResponse {
	return ProductResponse{
		ID:         product.ID.String(),
		TenantID:   product.TenantID.String(),
		Code:       product.Code,
		Name:       product.Name,
		Category:   product.Category,
		PriceCents: product.PriceCents,
		Status:     product.Status,
		IsActive:   product.IsActive,
	}
}
