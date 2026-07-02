package product

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

func (r Repository) List(ctx context.Context, tenantID uuid.UUID) ([]Product, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, tenant_id, code, name, category, price_cents, status, is_active
		FROM products
		WHERE tenant_id = $1
			AND is_deleted = false
		ORDER BY name ASC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]Product, 0)
	for rows.Next() {
		product := Product{}
		if err := rows.Scan(
			&product.ID,
			&product.TenantID,
			&product.Code,
			&product.Name,
			&product.Category,
			&product.PriceCents,
			&product.Status,
			&product.IsActive,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, rows.Err()
}

func (r Repository) Get(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (Product, error) {
	product := Product{}
	err := r.db.QueryRow(ctx, `
		SELECT id, tenant_id, code, name, category, price_cents, status, is_active
		FROM products
		WHERE tenant_id = $1
			AND id = $2
			AND is_deleted = false
	`, tenantID, id).Scan(
		&product.ID,
		&product.TenantID,
		&product.Code,
		&product.Name,
		&product.Category,
		&product.PriceCents,
		&product.Status,
		&product.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Product{}, ErrNotFound
		}
		return Product{}, err
	}
	return product, nil
}

func (r Repository) Create(ctx context.Context, product Product, actorID uuid.UUID, now time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO products (
			id, tenant_id, code, name, category, price_cents, status, is_active,
			created_by, created_date, updated_by, updated_date, is_deleted
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $9, $10, false)
	`, product.ID, product.TenantID, product.Code, product.Name, product.Category, product.PriceCents, product.Status, product.IsActive, actorID, now)
	return err
}

func (r Repository) Update(ctx context.Context, product Product, actorID uuid.UUID, now time.Time) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE products
		SET name = $3,
			category = $4,
			price_cents = $5,
			status = $6,
			is_active = $7,
			updated_by = $8,
			updated_date = $9
		WHERE tenant_id = $1
			AND id = $2
			AND is_deleted = false
	`, product.TenantID, product.ID, product.Name, product.Category, product.PriceCents, product.Status, product.IsActive, actorID, now)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r Repository) Delete(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID uuid.UUID, now time.Time) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE products
		SET is_deleted = true,
			is_active = false,
			updated_by = $3,
			updated_date = $4
		WHERE tenant_id = $1
			AND id = $2
			AND is_deleted = false
	`, tenantID, id, actorID, now)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
