package audit

import (
	"context"
	"encoding/json"
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

type Meta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

func NewService(repo Repository) Service {
	return Service{repo: repo, now: time.Now}
}

func (s Service) Write(ctx context.Context, event Event) error {
	return s.repo.Write(ctx, event, s.now().UTC())
}

func (s Service) List(ctx context.Context, query ListQuery) ([]LogResponse, Meta, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 20
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	logs, total, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, Meta{}, err
	}
	result := make([]LogResponse, 0, len(logs))
	for _, log := range logs {
		result = append(result, toResponse(log))
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + query.PageSize - 1) / query.PageSize
	}

	return result, Meta{
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

func (s Service) Get(ctx context.Context, id uuid.UUID) (LogResponse, error) {
	log, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return LogResponse{}, apperrors.New("NOT_FOUND", fiber.StatusNotFound, "Resource not found")
		}
		return LogResponse{}, err
	}
	return toResponse(log), nil
}

func toResponse(log Log) LogResponse {
	var tenantID *string
	if log.TenantID != nil {
		value := log.TenantID.String()
		tenantID = &value
	}
	var actorUserID *string
	if log.ActorUserID != nil {
		value := log.ActorUserID.String()
		actorUserID = &value
	}

	metadata := map[string]any{}
	if len(log.Metadata) > 0 {
		_ = json.Unmarshal(log.Metadata, &metadata)
	}

	return LogResponse{
		ID:           log.ID.String(),
		TenantID:     tenantID,
		ActorUserID:  actorUserID,
		Action:       log.Action,
		ResourceType: log.ResourceType,
		ResourceID:   log.ResourceID,
		IPAddress:    log.IPAddress,
		UserAgent:    log.UserAgent,
		Metadata:     metadata,
		CreatedDate:  log.CreatedDate,
	}
}
