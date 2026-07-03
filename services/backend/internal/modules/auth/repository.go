package auth

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

func (r Repository) FindUserByEmail(ctx context.Context, email string) (User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, tenant_id, email, name, password_hash, status, is_active, is_deleted
		FROM users
		WHERE lower(email) = lower($1)
			AND is_deleted = false
		LIMIT 1
	`, email)

	return scanUser(row)
}

func (r Repository) FindUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, tenant_id, email, name, password_hash, status, is_active, is_deleted
		FROM users
		WHERE id = $1
			AND is_deleted = false
		LIMIT 1
	`, id)

	return scanUser(row)
}

func (r Repository) UpdateLastLogin(ctx context.Context, userID uuid.UUID, now time.Time) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users
		SET last_login_at = $2,
			updated_by = $1,
			updated_date = $2
		WHERE id = $1
			AND is_deleted = false
	`, userID, now)
	return err
}

func (r Repository) CreateRefreshToken(ctx context.Context, token RefreshToken, actorID uuid.UUID, now time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO refresh_tokens (
			id, user_id, token_hash, expires_at, is_active,
			created_by, created_date, updated_by, updated_date, is_deleted
		) VALUES ($1, $2, $3, $4, true, $5, $6, $5, $6, false)
	`, token.ID, token.UserID, token.TokenHash, token.ExpiresAt, actorID, now)
	return err
}

func (r Repository) FindRefreshTokenByHash(ctx context.Context, tokenHash string) (RefreshToken, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, user_id, token_hash, expires_at, revoked_at
		FROM refresh_tokens
		WHERE token_hash = $1
			AND is_deleted = false
		LIMIT 1
	`, tokenHash)

	token := RefreshToken{}
	if err := row.Scan(&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.RevokedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return RefreshToken{}, ErrNotFound
		}
		return RefreshToken{}, err
	}
	return token, nil
}

func (r Repository) RevokeRefreshToken(ctx context.Context, tokenID uuid.UUID, actorID uuid.UUID, replacedBy *uuid.UUID, now time.Time) error {
	_, err := r.db.Exec(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = $2,
			replaced_by_token_id = $3,
			is_active = false,
			updated_by = $4,
			updated_date = $2
		WHERE id = $1
			AND is_deleted = false
	`, tokenID, now, replacedBy, actorID)
	return err
}

func (r Repository) EnsureTenant(ctx context.Context, id uuid.UUID, name string, slug string, now time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO tenants (
			id, name, slug, status, is_active,
			created_by, created_date, updated_by, updated_date, is_deleted
		) VALUES ($1, $2, $3, 'active', true, NULL, $4, NULL, $4, false)
		ON CONFLICT (slug)
		DO UPDATE SET name = $2,
			status = 'active',
			is_active = true,
			is_deleted = false,
			updated_date = $4
	`, id, name, slug, now)
	return err
}

func (r Repository) FindTenantBySlug(ctx context.Context, slug string) (uuid.UUID, error) {
	var id uuid.UUID
	if err := r.db.QueryRow(ctx, `
		SELECT id
		FROM tenants
		WHERE slug = $1
			AND is_deleted = false
		LIMIT 1
	`, slug).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, ErrNotFound
		}
		return uuid.Nil, err
	}
	return id, nil
}

func (r Repository) EnsureUser(ctx context.Context, user User, actorID uuid.UUID, now time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO users (
			id, tenant_id, email, name, password_hash, status, is_active,
			created_by, created_date, updated_by, updated_date, is_deleted
		) VALUES ($1, $2, $3, $4, $5, $6, true, $7, $8, $7, $8, false)
		ON CONFLICT DO NOTHING
	`, user.ID, user.TenantID, user.Email, user.Name, user.PasswordHash, user.Status, actorID, now)
	return err
}

func (r Repository) FindUserByEmailAndTenant(ctx context.Context, email string, tenantID uuid.UUID) (User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, tenant_id, email, name, password_hash, status, is_active, is_deleted
		FROM users
		WHERE lower(email) = lower($1)
			AND tenant_id = $2
			AND is_deleted = false
		LIMIT 1
	`, email, tenantID)

	return scanUser(row)
}

func scanUser(row pgx.Row) (User, error) {
	user := User{}
	if err := row.Scan(
		&user.ID,
		&user.TenantID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.Status,
		&user.IsActive,
		&user.IsDeleted,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}
	return user, nil
}
