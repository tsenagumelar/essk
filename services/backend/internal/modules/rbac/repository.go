package rbac

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

const userHasPermissionSQL = `
	SELECT EXISTS (
		SELECT 1
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles ro ON ro.id = ur.role_id
		JOIN role_permissions rp ON rp.role_id = ro.id
		JOIN permissions pe ON pe.id = rp.permission_id
		WHERE u.id = $1
			AND pe.code = $2
			AND u.is_active = true
			AND u.is_deleted = false
			AND ur.is_active = true
			AND ur.is_deleted = false
			AND ro.is_active = true
			AND ro.is_deleted = false
			AND (ro.tenant_id IS NULL OR ro.tenant_id = u.tenant_id)
			AND rp.is_active = true
			AND rp.is_deleted = false
			AND pe.is_active = true
			AND pe.is_deleted = false
	)
`

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return Repository{db: db}
}

func (r Repository) UserHasPermission(ctx context.Context, userID uuid.UUID, permissionCode string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, userHasPermissionSQL, userID, permissionCode).Scan(&exists)
	return exists, err
}

func RoleAppliesToUserTenant(userTenantID *uuid.UUID, roleTenantID *uuid.UUID) bool {
	if roleTenantID == nil {
		return true
	}
	if userTenantID == nil {
		return false
	}
	return *roleTenantID == *userTenantID
}

func (r Repository) ListPermissions(ctx context.Context) ([]Permission, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, code, name, description, is_active
		FROM permissions
		WHERE is_deleted = false
		ORDER BY code ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := make([]Permission, 0)
	for rows.Next() {
		permission := Permission{}
		if err := rows.Scan(&permission.ID, &permission.Code, &permission.Name, &permission.Description, &permission.IsActive); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, rows.Err()
}

func (r Repository) ListRoles(ctx context.Context, tenantID *uuid.UUID) ([]Role, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, tenant_id, name, code, description, is_system, is_active
		FROM roles
		WHERE is_deleted = false
			AND ($1::uuid IS NULL OR tenant_id = $1)
		ORDER BY code ASC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]Role, 0)
	for rows.Next() {
		role := Role{}
		if err := rows.Scan(&role.ID, &role.TenantID, &role.Name, &role.Code, &role.Description, &role.IsSystem, &role.IsActive); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (r Repository) GetRole(ctx context.Context, id uuid.UUID) (Role, error) {
	role := Role{}
	err := r.db.QueryRow(ctx, `
		SELECT id, tenant_id, name, code, description, is_system, is_active
		FROM roles
		WHERE id = $1
			AND is_deleted = false
	`, id).Scan(&role.ID, &role.TenantID, &role.Name, &role.Code, &role.Description, &role.IsSystem, &role.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Role{}, ErrNotFound
		}
		return Role{}, err
	}
	return role, nil
}

func (r Repository) CreateRole(ctx context.Context, role Role, actorID uuid.UUID, now time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO roles (
			id, tenant_id, name, code, description, is_system, is_active,
			created_by, created_date, updated_by, updated_date, is_deleted
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $8, $9, false)
	`, role.ID, role.TenantID, role.Name, role.Code, role.Description, role.IsSystem, role.IsActive, actorID, now)
	return err
}

func (r Repository) UpdateRole(ctx context.Context, role Role, actorID uuid.UUID, now time.Time) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE roles
		SET name = $2,
			description = $3,
			is_active = $4,
			updated_by = $5,
			updated_date = $6
		WHERE id = $1
			AND is_deleted = false
	`, role.ID, role.Name, role.Description, role.IsActive, actorID, now)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r Repository) DeleteRole(ctx context.Context, id uuid.UUID, actorID uuid.UUID, now time.Time) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE roles
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

func (r Repository) AssignPermission(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID, actorID uuid.UUID, now time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO role_permissions (
			role_id, permission_id, is_active,
			created_by, created_date, updated_by, updated_date, is_deleted
		) VALUES ($1, $2, true, $3, $4, $3, $4, false)
		ON CONFLICT (role_id, permission_id)
		DO UPDATE SET is_active = true, is_deleted = false, updated_by = $3, updated_date = $4
	`, roleID, permissionID, actorID, now)
	return err
}

func (r Repository) RemovePermission(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID, actorID uuid.UUID, now time.Time) error {
	_, err := r.db.Exec(ctx, `
		UPDATE role_permissions
		SET is_active = false,
			is_deleted = true,
			updated_by = $3,
			updated_date = $4
		WHERE role_id = $1
			AND permission_id = $2
	`, roleID, permissionID, actorID, now)
	return err
}

func (r Repository) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID, actorID uuid.UUID, now time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO user_roles (
			user_id, role_id, is_active,
			created_by, created_date, updated_by, updated_date, is_deleted
		) VALUES ($1, $2, true, $3, $4, $3, $4, false)
		ON CONFLICT (user_id, role_id)
		DO UPDATE SET is_active = true, is_deleted = false, updated_by = $3, updated_date = $4
	`, userID, roleID, actorID, now)
	return err
}

func (r Repository) RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID, actorID uuid.UUID, now time.Time) error {
	_, err := r.db.Exec(ctx, `
		UPDATE user_roles
		SET is_active = false,
			is_deleted = true,
			updated_by = $3,
			updated_date = $4
		WHERE user_id = $1
			AND role_id = $2
	`, userID, roleID, actorID, now)
	return err
}

func (r Repository) EnsurePermission(ctx context.Context, permission Permission, actorID *uuid.UUID, now time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO permissions (
			id, code, name, description, is_active,
			created_by, created_date, updated_by, updated_date, is_deleted
		) VALUES ($1, $2, $3, $4, true, $5, $6, $5, $6, false)
		ON CONFLICT (code)
		DO UPDATE SET name = $3, description = $4, is_active = true, updated_by = $5, updated_date = $6, is_deleted = false
	`, permission.ID, permission.Code, permission.Name, permission.Description, actorID, now)
	return err
}

func (r Repository) EnsureRole(ctx context.Context, role Role, actorID uuid.UUID, now time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO roles (
			id, tenant_id, name, code, description, is_system, is_active,
			created_by, created_date, updated_by, updated_date, is_deleted
		) VALUES ($1, $2, $3, $4, $5, $6, true, $7, $8, $7, $8, false)
		ON CONFLICT DO NOTHING
	`, role.ID, role.TenantID, role.Name, role.Code, role.Description, role.IsSystem, actorID, now)
	return err
}

func (r Repository) FindRoleByCode(ctx context.Context, tenantID *uuid.UUID, code string) (Role, error) {
	role := Role{}
	err := r.db.QueryRow(ctx, `
		SELECT id, tenant_id, name, code, description, is_system, is_active
		FROM roles
		WHERE code = $2
			AND is_deleted = false
			AND (($1::uuid IS NULL AND tenant_id IS NULL) OR tenant_id = $1)
		LIMIT 1
	`, tenantID, code).Scan(&role.ID, &role.TenantID, &role.Name, &role.Code, &role.Description, &role.IsSystem, &role.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Role{}, ErrNotFound
		}
		return Role{}, err
	}
	return role, nil
}

func (r Repository) FindPermissionByCode(ctx context.Context, code string) (Permission, error) {
	permission := Permission{}
	err := r.db.QueryRow(ctx, `
		SELECT id, code, name, description, is_active
		FROM permissions
		WHERE code = $1
			AND is_deleted = false
		LIMIT 1
	`, code).Scan(&permission.ID, &permission.Code, &permission.Name, &permission.Description, &permission.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Permission{}, ErrNotFound
		}
		return Permission{}, err
	}
	return permission, nil
}
