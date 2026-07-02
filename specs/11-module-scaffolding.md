# ESSK Specs - Module Scaffolding

## Purpose

ESSK must include a scaffolding command so new modules are generated consistently across backend, frontend, API, database, and tests.

Primary command:

```text
essk add-module user-management
```

## CLI Location

Recommended implementation:

```text
services/backend/cmd/essk/
```

The CLI is written in Go so it can reuse backend naming helpers and templates.

Recommended library:

- `spf13/cobra` for CLI commands.

Template engine:

- Go `text/template`.

Template location:

```text
tools/scaffold/templates/
  backend/
  frontend/
  migrations/
  queries/
  docs/
```

## Command Options

```text
essk add-module {name}
  --tenant-scoped=true
  --with-backend=true
  --with-frontend=true
  --with-admin=false
  --with-tests=true
  --force=false
```

Default behavior:

- Backend: enabled.
- Frontend web CRUD: enabled.
- Admin app: disabled until `apps/admin` is bootstrapped.
- Tests: enabled.
- Tenant-scoped: enabled for business modules.

## Naming Rules

Input:

```text
user-management
```

Derived names:

- Route path: `user-management`.
- Go package: `user_management`.
- Database table: `user_management`.
- Type name: `UserManagement`.
- Variable name: `userManagement`.
- Permission prefix: `user_management`.

Generated permissions:

```text
user_management:read
user_management:create
user_management:update
user_management:delete
```

## Backend Generated Files

```text
services/backend/internal/modules/user_management/
  entity.go
  dto.go
  commands.go
  queries.go
  command_handler.go
  query_handler.go
  repository.go
  service.go
  handler.go
  route.go
  validation.go
  errors.go
  audit.go
  command_handler_test.go
  query_handler_test.go
  handler_test.go

services/backend/queries/user_management.sql
services/backend/migrations/{version}_create_user_management.up.sql
services/backend/migrations/{version}_create_user_management.down.sql
```

Generated endpoints:

```text
GET    /api/v1/user-management
POST   /api/v1/user-management
GET    /api/v1/user-management/:id
PATCH  /api/v1/user-management/:id
DELETE /api/v1/user-management/:id
```

Generated backend behavior:

- List with pagination.
- Get by ID.
- Create.
- Update.
- Soft delete through `is_deleted = true`.
- Active/inactive support through `is_active`.
- Tenant filtering when `--tenant-scoped=true`.
- CQRS command/query separation.
- Permission middleware on all CRUD routes.
- Audit log calls for create, update, and delete.

## Migration Template

Every generated table must include:

```sql
id uuid primary key,
is_active boolean not null default true,
created_by uuid null,
created_date timestamptz not null default now(),
updated_by uuid null,
updated_date timestamptz not null default now(),
is_deleted boolean not null default false
```

For tenant-scoped modules:

```sql
tenant_id uuid null references tenants(id)
```

Default generated business fields:

```sql
name varchar(160) not null,
description text null
```

Default generated indexes:

```sql
create index idx_{table}_is_deleted on {table}(is_deleted);
create index idx_{table}_is_active on {table}(is_active);
create index idx_{table}_created_date on {table}(created_date);
```

Tenant-scoped modules also require:

```sql
create index idx_{table}_tenant_id on {table}(tenant_id);
```

## Frontend Generated Files

```text
apps/web/app/(app)/user-management/
  page.tsx
  loading.tsx
  error.tsx
  new/
    page.tsx
  [id]/
    page.tsx
    edit/
      page.tsx

apps/web/features/user-management/
  api.ts
  schema.ts
  types.ts
  hooks.ts
  components/
    user-management-form.tsx
    user-management-table.tsx
    user-management-delete-dialog.tsx
```

Generated UI behavior:

- List page.
- Search input.
- Pagination.
- Create button.
- Edit action.
- Delete action with confirmation.
- Detail view.
- Create and edit forms.
- Loading state.
- Error state.
- Empty state.

Generated frontend data layer:

- `listUserManagement`.
- `getUserManagement`.
- `createUserManagement`.
- `updateUserManagement`.
- `deleteUserManagement`.
- TanStack Query hooks for each operation.
- Zod schema for create and update.

## Registry Updates

The scaffolder should update central registries when they exist:

- Backend route registry.
- Backend permission seed list.
- Frontend navigation config.
- Frontend route metadata.

If a registry does not exist yet, generated files must include a short TODO pointing to the expected manual registration location.

## Safety Rules

- Fail when target files exist unless `--force=true`.
- Print a dry-run file list before writing when `--dry-run=true`.
- Never delete existing files.
- Never overwrite unrelated modules.
- Format generated Go files with `gofmt`.
- Format generated TypeScript files with project formatter when available.

## Verification Command

After generation, developer should be able to run:

```text
go test ./...
pnpm lint
pnpm typecheck
```

When OpenAPI generation is available, scaffolding should also run or instruct:

```text
pnpm codegen
```
