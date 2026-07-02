package tenant

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
	_ = s.writeAudit(ctx, actorID, "tenant.create", tenant.ID.String(), map[string]any{"slug": tenant.Slug})
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
	_ = s.writeAudit(ctx, actorID, "tenant.update", tenant.ID.String(), map[string]any{"status": tenant.Status})
	return toResponse(tenant), nil
}

func (s Service) Delete(ctx context.Context, id uuid.UUID, actorID uuid.UUID) error {
	if err := mapNotFound(s.repo.Delete(ctx, id, actorID, s.now().UTC())); err != nil {
		return err
	}
	_ = s.writeAudit(ctx, actorID, "tenant.delete", id.String(), nil)
	return nil
}

func (s Service) writeAudit(ctx context.Context, actorID uuid.UUID, action string, resourceID string, metadata map[string]any) error {
	if s.audit == nil {
		return nil
	}
	return s.audit.Write(ctx, audit.Event{
		ActorUserID:  &actorID,
		Action:       action,
		ResourceType: "tenant",
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

func toResponse(tenant Tenant) TenantResponse {
	return TenantResponse{
		ID:       tenant.ID.String(),
		Name:     tenant.Name,
		Slug:     tenant.Slug,
		Status:   tenant.Status,
		IsActive: tenant.IsActive,
	}
}
