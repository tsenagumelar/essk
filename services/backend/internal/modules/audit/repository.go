package audit

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return Repository{db: db}
}

func (r Repository) Write(ctx context.Context, event Event, now time.Time) error {
	metadata := []byte("{}")
	if event.Metadata != nil {
		encoded, err := json.Marshal(event.Metadata)
		if err != nil {
			return err
		}
		metadata = encoded
	}

	var createdBy *uuid.UUID
	if event.ActorUserID != nil {
		createdBy = event.ActorUserID
	}

	_, err := r.db.Exec(ctx, `
		INSERT INTO audit_logs (
			id, tenant_id, actor_user_id, action, resource_type, resource_id,
			ip_address, user_agent, metadata, is_active,
			created_by, created_date, updated_by, updated_date, is_deleted
		) VALUES ($1, $2, $3, $4, $5, $6, $7::inet, $8, $9, true, $10, $11, $10, $11, false)
	`, uuid.New(), event.TenantID, event.ActorUserID, event.Action, event.ResourceType, event.ResourceID, event.IPAddress, event.UserAgent, metadata, createdBy, now)
	return err
}

func (r Repository) List(ctx context.Context, query ListQuery) ([]Log, int, error) {
	page := query.Page
	if page < 1 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var tenantID *uuid.UUID
	if query.TenantID != nil && *query.TenantID != "" {
		parsed, err := uuid.Parse(*query.TenantID)
		if err != nil {
			return nil, 0, err
		}
		tenantID = &parsed
	}
	var actorUserID *uuid.UUID
	if query.ActorUserID != nil && *query.ActorUserID != "" {
		parsed, err := uuid.Parse(*query.ActorUserID)
		if err != nil {
			return nil, 0, err
		}
		actorUserID = &parsed
	}

	var total int
	if err := r.db.QueryRow(ctx, `
		SELECT count(*)
		FROM audit_logs
		WHERE is_deleted = false
			AND ($1::uuid IS NULL OR tenant_id = $1)
			AND ($2::uuid IS NULL OR actor_user_id = $2)
			AND ($3 = '' OR action = $3)
			AND ($4 = '' OR resource_type = $4)
	`, tenantID, actorUserID, query.Action, query.ResourceType).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, tenant_id, actor_user_id, action, resource_type, resource_id,
			host(ip_address), user_agent, metadata, created_date
		FROM audit_logs
		WHERE is_deleted = false
			AND ($1::uuid IS NULL OR tenant_id = $1)
			AND ($2::uuid IS NULL OR actor_user_id = $2)
			AND ($3 = '' OR action = $3)
			AND ($4 = '' OR resource_type = $4)
		ORDER BY created_date DESC
		LIMIT $5 OFFSET $6
	`, tenantID, actorUserID, query.Action, query.ResourceType, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	logs := make([]Log, 0)
	for rows.Next() {
		log := Log{}
		if err := rows.Scan(
			&log.ID,
			&log.TenantID,
			&log.ActorUserID,
			&log.Action,
			&log.ResourceType,
			&log.ResourceID,
			&log.IPAddress,
			&log.UserAgent,
			&log.Metadata,
			&log.CreatedDate,
		); err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	return logs, total, rows.Err()
}

func (r Repository) Get(ctx context.Context, id uuid.UUID) (Log, error) {
	log := Log{}
	err := r.db.QueryRow(ctx, `
		SELECT id, tenant_id, actor_user_id, action, resource_type, resource_id,
			host(ip_address), user_agent, metadata, created_date
		FROM audit_logs
		WHERE id = $1
			AND is_deleted = false
	`, id).Scan(
		&log.ID,
		&log.TenantID,
		&log.ActorUserID,
		&log.Action,
		&log.ResourceType,
		&log.ResourceID,
		&log.IPAddress,
		&log.UserAgent,
		&log.Metadata,
		&log.CreatedDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Log{}, ErrNotFound
		}
		return Log{}, err
	}
	return log, nil
}
