# ESSK Specs - API, Auth, And RBAC

## API Standard

Base URL:

```text
/api/v1
```

Success response:

```json
{
  "success": true,
  "message": "OK",
  "data": {},
  "meta": {}
}
```

Error response:

```json
{
  "success": false,
  "message": "Validation Error",
  "errors": [
    {
      "field": "email",
      "code": "required",
      "message": "Email is required"
    }
  ]
}
```

## Endpoint Groups

Authentication:

```text
POST /auth/login
POST /auth/logout
POST /auth/refresh
GET  /auth/me
```

Users:

```text
GET    /users
POST   /users
GET    /users/:id
PATCH  /users/:id
DELETE /users/:id
GET    /users/me
PATCH  /users/me
```

Tenants:

```text
GET    /tenants
POST   /tenants
GET    /tenants/:id
PATCH  /tenants/:id
DELETE /tenants/:id
```

Roles:

```text
GET    /roles
POST   /roles
GET    /roles/:id
PATCH  /roles/:id
DELETE /roles/:id
POST   /roles/:id/permissions
DELETE /roles/:id/permissions/:permission_id
```

Permissions:

```text
GET /permissions
```

Audit:

```text
GET /audit-logs
GET /audit-logs/:id
```

## Authentication

Token model:

- Access token: JWT, short TTL, default 15 minutes.
- Refresh token: opaque random token, stored hashed in database, default TTL 7 days.

JWT claims:

- `sub`: user ID.
- `tenant_id`: tenant ID.
- `roles`: role codes.
- `permissions`: permission codes or a permissions version reference.
- `iss`: configured issuer.
- `iat`: issued at.
- `exp`: expires at.

Password hashing:

- Use argon2id.
- Store encoded hash string with algorithm parameters.
- Never log password or token values.

Login behavior:

1. Validate email and password.
2. Find active user.
3. Verify password hash.
4. Create access token.
5. Create refresh token record with token hash.
6. Write audit log.
7. Return tokens and profile.

Logout behavior:

- Revoke current refresh token.
- Optionally denylist current access token until expiry if Redis is enabled.
- Write audit log.

Refresh behavior:

- Verify refresh token hash.
- Ensure token is not expired or revoked.
- Rotate refresh token.
- Return new access and refresh tokens.

## Authorization

Use permission-based authorization.

Role is a grouping mechanism. API checks should use permission code, not role name.

Permission naming:

```text
resource:action
```

Examples:

```text
users:read
users:create
users:update
users:delete
tenants:read
roles:manage
audit_logs:read
```

Middleware:

```text
RequireAuth()
RequirePermission("users:read")
```

## Tenant Boundary

Every authenticated request should carry tenant context when applicable.

Tenant source priority:

1. JWT claim.
2. Explicit tenant header only for system/admin flows.
3. Route param only when endpoint is designed for cross-tenant administration.

Default rule:

- A tenant user can only access resources within their tenant.
- System admin endpoints must be explicit and permission-protected.

## Pagination

Query params:

```text
page=1
page_size=20
sort=created_date:desc
search=keyword
```

Limits:

- Default page size: 20.
- Max page size: 100.

## OpenAPI

Backend must expose:

```text
/swagger/index.html
/swagger/doc.json
```

OpenAPI must describe:

- All endpoints.
- Request DTOs.
- Response DTOs.
- Error responses.
- Auth scheme.
- Pagination params.

Frontend SDK generation depends on this contract.
