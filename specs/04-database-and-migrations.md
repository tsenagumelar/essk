# ESSK Specs - Database And Migrations

## Database

Primary database: PostgreSQL 16+.

Use PostgreSQL as the source of truth for:

- Users.
- Tenants.
- Roles.
- Permissions.
- Audit logs.
- Refresh token records.
- Application metadata.

Redis is used for:

- Short-lived cache.
- Optional token denylist.
- Rate limiting store.
- Future background job backend.

## Query Strategy

Use sqlc + pgx.

Reasoning:

- SQL remains explicit.
- Generated Go types reduce runtime mapping bugs.
- Easy to optimize queries.
- Better fit for enterprise reporting and audit queries than implicit ORM behavior.

sqlc config:

```text
services/backend/sqlc.yaml
```

Recommended generated package:

```text
internal/database/sqlc
```

## Migration Tool

Use `golang-migrate/migrate`.

Migration location:

```text
services/backend/migrations/
```

Naming:

```text
000001_create_tenants.up.sql
000001_create_tenants.down.sql
000002_create_users.up.sql
000002_create_users.down.sql
```

Rules:

- Every schema change requires up and down migration.
- No destructive migration without explicit review.
- Avoid editing committed migrations.
- Seed data should be separate from schema migration unless required for app boot.

## Mandatory Columns For All Tables

Every table must include these standard lifecycle columns:

- `is_active boolean not null default true`.
- `created_by uuid null`.
- `created_date timestamptz not null default now()`.
- `updated_by uuid null`.
- `updated_date timestamptz not null default now()`.
- `is_deleted boolean not null default false`.

Rules:

- Use `created_date` and `updated_date`, not `created_at` or `updated_at`.
- Use soft delete through `is_deleted = true`; physical delete is only allowed for temporary tables or explicit maintenance jobs.
- Default read queries must filter `is_deleted = false`.
- Default active queries should filter `is_active = true` when the table represents an enabled/disabled business object.
- `created_by` and `updated_by` reference `users.id` conceptually, but the foreign key can be omitted for low-level system tables to avoid circular bootstrap issues.
- Join tables also include the mandatory columns unless explicitly documented as a pure internal temporary table.
- Service layer must set `created_by` and `updated_by` from authenticated actor context when available.

## Core Tables

### `tenants`

Purpose: tenant boundary for multi-tenant SaaS.

Columns:

- `id uuid primary key`.
- `name varchar(160) not null`.
- `slug varchar(120) not null unique`.
- `status varchar(32) not null`.
- Mandatory lifecycle columns.

### `users`

Columns:

- `id uuid primary key`.
- `tenant_id uuid null references tenants(id)`.
- `email varchar(255) not null`.
- `name varchar(160) not null`.
- `password_hash text not null`.
- `status varchar(32) not null`.
- `last_login_at timestamptz null`.
- Mandatory lifecycle columns.

Unique constraints:

- `unique(tenant_id, email)` for tenant-scoped users.

### `roles`

Columns:

- `id uuid primary key`.
- `tenant_id uuid null references tenants(id)`.
- `name varchar(120) not null`.
- `code varchar(120) not null`.
- `description text null`.
- `is_system boolean not null default false`.
- Mandatory lifecycle columns.

### `permissions`

Columns:

- `id uuid primary key`.
- `code varchar(160) not null unique`.
- `name varchar(160) not null`.
- `description text null`.
- Mandatory lifecycle columns.

### `role_permissions`

Columns:

- `role_id uuid not null references roles(id)`.
- `permission_id uuid not null references permissions(id)`.
- Mandatory lifecycle columns.

Primary key:

- `(role_id, permission_id)`.

### `user_roles`

Columns:

- `user_id uuid not null references users(id)`.
- `role_id uuid not null references roles(id)`.
- Mandatory lifecycle columns.

Primary key:

- `(user_id, role_id)`.

### `refresh_tokens`

Columns:

- `id uuid primary key`.
- `user_id uuid not null references users(id)`.
- `token_hash text not null`.
- `expires_at timestamptz not null`.
- `revoked_at timestamptz null`.
- Mandatory lifecycle columns.

### `audit_logs`

Columns:

- `id uuid primary key`.
- `tenant_id uuid null`.
- `actor_user_id uuid null`.
- `action varchar(120) not null`.
- `resource_type varchar(120) not null`.
- `resource_id varchar(120) null`.
- `ip_address inet null`.
- `user_agent text null`.
- `metadata jsonb not null default '{}'::jsonb`.
- Mandatory lifecycle columns.

## Multi-Tenancy Strategy

Initial strategy: shared database, shared schema, tenant discriminator column.

Rules:

- Tenant-scoped tables include `tenant_id`.
- Service layer must enforce tenant boundary.
- Repository methods should accept tenant ID for tenant-scoped queries.
- Super admin/system records may use `tenant_id null` only when explicitly designed.

Future strategy:

- Schema-per-tenant or database-per-tenant can be introduced later for selected enterprise customers.

## Seed Data

Foundation seed:

- Default system tenant.
- Admin role.
- Basic permissions.
- Default admin user.

Seed credentials must be configurable through local-only environment variables and must not be hardcoded for production.
