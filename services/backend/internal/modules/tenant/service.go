package tenant

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	apperrors "github.com/tsenagumelar/essk/services/backend/internal/errors"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) Service {
	return Service{repo: repo, now: time.Now}
}

func (s Service) List(ctx context.Context) ([]TenantResponse, error) {
	tenants, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]TenantResponse, 0, len(tenants))
	for _, tenant := range tenants {
		result = append(result, toResponse(tenant))
	}
	return result, nil
}

func (s Service) Get(ctx context.Context, id uuid.UUID) (TenantResponse, error) {
	tenant, err := s.repo.Get(ctx, id)
	if err != nil {
		return TenantResponse{}, mapNotFound(err)
	}
	return toResponse(tenant), nil
}

func (s Service) Create(ctx context.Context, req CreateTenantRequest, actorID uuid.UUID) (TenantResponse, error) {
	tenant := Tenant{
		ID:       uuid.New(),
		Name:     req.Name,
		Slug:     req.Slug,
		Status:   "active",
		IsActive: true,
	}
	if err := s.repo.Create(ctx, tenant, actorID, s.now().UTC()); err != nil {
		return TenantResponse{}, err
	}
	return toResponse(tenant), nil
}

func (s Service) Update(ctx context.Context, id uuid.UUID, req UpdateTenantRequest, actorID uuid.UUID) (TenantResponse, error) {
	tenant, err := s.repo.Get(ctx, id)
	if err != nil {
		return TenantResponse{}, mapNotFound(err)
	}
	tenant.Name = req.Name
	tenant.Status = req.Status
	tenant.IsActive = req.IsActive
	if err := s.repo.Update(ctx, tenant, actorID, s.now().UTC()); err != nil {
		return TenantResponse{}, mapNotFound(err)
	}
	return toResponse(tenant), nil
}

func (s Service) Delete(ctx context.Context, id uuid.UUID, actorID uuid.UUID) error {
	return mapNotFound(s.repo.Delete(ctx, id, actorID, s.now().UTC()))
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

func toResponse(tenant Tenant) TenantResponse {
	return TenantResponse{
		ID:       tenant.ID.String(),
		Name:     tenant.Name,
		Slug:     tenant.Slug,
		Status:   tenant.Status,
		IsActive: tenant.IsActive,
	}
}
