package useradmin

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

func (r Repository) UserHasRoleCode(ctx context.Context, userID uuid.UUID, code string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM user_roles ur
			JOIN roles ro ON ro.id = ur.role_id
			WHERE ur.user_id = $1
				AND ro.code = $2
				AND ur.is_active = true
				AND ur.is_deleted = false
				AND ro.is_active = true
				AND ro.is_deleted = false
		)
	`, userID, code).Scan(&exists)
	return exists, err
}

func (r Repository) List(ctx context.Context, tenantID *uuid.UUID) ([]User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, tenant_id, email, name, status, is_active
		FROM users
		WHERE is_deleted = false
			AND ($1::uuid IS NULL OR tenant_id = $1)
		ORDER BY name ASC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.TenantID, &user.Email, &user.Name, &user.Status, &user.IsActive); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r Repository) Get(ctx context.Context, id uuid.UUID, tenantID *uuid.UUID) (User, error) {
	user := User{}
	err := r.db.QueryRow(ctx, `
		SELECT id, tenant_id, email, name, status, is_active
		FROM users
		WHERE id = $1
			AND is_deleted = false
			AND ($2::uuid IS NULL OR tenant_id = $2)
	`, id, tenantID).Scan(&user.ID, &user.TenantID, &user.Email, &user.Name, &user.Status, &user.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}
	return user, nil
}

func (r Repository) Create(ctx context.Context, user User, passwordHash string, actorID uuid.UUID, now time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO users (
			id, tenant_id, email, name, password_hash, status, is_active,
			created_by, created_date, updated_by, updated_date, is_deleted
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $8, $9, false)
	`, user.ID, user.TenantID, user.Email, user.Name, passwordHash, user.Status, user.IsActive, actorID, now)
	return err
}

func (r Repository) Update(ctx context.Context, user User, actorID uuid.UUID, now time.Time) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE users
		SET name = $2,
			status = $3,
			is_active = $4,
			updated_by = $5,
			updated_date = $6
		WHERE id = $1
			AND is_deleted = false
	`, user.ID, user.Name, user.Status, user.IsActive, actorID, now)
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
		UPDATE users
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

func (r Repository) ListRoleIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, `
		SELECT role_id
		FROM user_roles
		WHERE user_id = $1
			AND is_active = true
			AND is_deleted = false
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roleIDs := make([]uuid.UUID, 0)
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		roleIDs = append(roleIDs, id)
	}
	return roleIDs, rows.Err()
}

func (r Repository) ReplaceRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID, actorID uuid.UUID, now time.Time) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		UPDATE user_roles
		SET is_active = false,
			is_deleted = true,
			updated_by = $2,
			updated_date = $3
		WHERE user_id = $1
	`, userID, actorID, now); err != nil {
		return err
	}

	for _, roleID := range roleIDs {
		if _, err := tx.Exec(ctx, `
			INSERT INTO user_roles (
				user_id, role_id, is_active,
				created_by, created_date, updated_by, updated_date, is_deleted
			) VALUES ($1, $2, true, $3, $4, $3, $4, false)
			ON CONFLICT (user_id, role_id)
			DO UPDATE SET is_active = true, is_deleted = false, updated_by = $3, updated_date = $4
		`, userID, roleID, actorID, now); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
