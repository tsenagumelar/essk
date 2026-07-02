package tenant

import (
	"context"
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

func (r Repository) List(ctx context.Context) ([]Tenant, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, slug, status, is_active
		FROM tenants
		WHERE is_deleted = false
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenants := make([]Tenant, 0)
	for rows.Next() {
		tenant := Tenant{}
		if err := rows.Scan(&tenant.ID, &tenant.Name, &tenant.Slug, &tenant.Status, &tenant.IsActive); err != nil {
			return nil, err
		}
		tenants = append(tenants, tenant)
	}
	return tenants, rows.Err()
}

func (r Repository) Get(ctx context.Context, id uuid.UUID) (Tenant, error) {
	tenant := Tenant{}
	err := r.db.QueryRow(ctx, `
		SELECT id, name, slug, status, is_active
		FROM tenants
		WHERE id = $1
			AND is_deleted = false
	`, id).Scan(&tenant.ID, &tenant.Name, &tenant.Slug, &tenant.Status, &tenant.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Tenant{}, ErrNotFound
		}
		return Tenant{}, err
	}
	return tenant, nil
}

func (r Repository) Create(ctx context.Context, tenant Tenant, actorID uuid.UUID, now time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO tenants (
			id, name, slug, status, is_active,
			created_by, created_date, updated_by, updated_date, is_deleted
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $6, $7, false)
	`, tenant.ID, tenant.Name, tenant.Slug, tenant.Status, tenant.IsActive, actorID, now)
	return err
}

func (r Repository) Update(ctx context.Context, tenant Tenant, actorID uuid.UUID, now time.Time) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE tenants
		SET name = $2,
			status = $3,
			is_active = $4,
			updated_by = $5,
			updated_date = $6
		WHERE id = $1
			AND is_deleted = false
	`, tenant.ID, tenant.Name, tenant.Status, tenant.IsActive, actorID, now)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r Repository) Delete(ctx context.Context, id uuid.UUID, actorID uuid.UUID, now time.Time) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE tenants
		SET is_deleted = true,
			is_active = false,
			updated_by = $2,
			updated_date = $3
		WHERE id = $1
			AND is_deleted = false
	`, id, actorID, now)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
