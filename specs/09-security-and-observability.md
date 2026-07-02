# ESSK Specs - Security And Observability

## Security Baseline

Mandatory:

- JWT access token.
- Refresh token rotation.
- Password hashing with argon2id.
- CORS allowlist.
- Security headers.
- Audit log.
- Request ID.
- Correlation ID.
- Input validation.
- Tenant boundary enforcement.
- No secrets in logs.
- Rate limiting for public and sensitive endpoints.
- OWASP Top 10 controls as acceptance criteria.

## HTTP Security Headers

Apply through backend and/or Nginx:

- `X-Content-Type-Options: nosniff`.
- `X-Frame-Options: DENY`.
- `Referrer-Policy: no-referrer`.
- `Content-Security-Policy` for frontend when stable.
- `Strict-Transport-Security` in HTTPS environments.

## CORS

Default local origins:

```text
http://localhost:3000
http://localhost:3001
```

Production origins must be explicit.

Do not use wildcard origins with credentials.

## Rate Limiting

Rate limiting is mandatory for enterprise readiness.

Implementation:

- Use Redis-backed distributed limiter for multi-instance deployments.
- Use route-specific limits instead of one global limit.
- Return standard API error response with HTTP `429`.
- Include `Retry-After` header when possible.
- Log rate-limit events with request ID, IP, user ID, tenant ID, and route.
- Exclude health endpoints from strict limits.

Recommended library:

- Fiber limiter middleware with a Redis storage adapter.

Default limits:

- Global API: 300 requests per minute per IP.
- Auth login: 5 requests per minute per IP and email pair.
- Auth refresh: 20 requests per minute per user/session.
- Write endpoints: 60 requests per minute per user.
- Admin endpoints: 120 requests per minute per user.

Priority endpoints:

- `/auth/login`.
- `/auth/refresh`.
- Password reset endpoints when added.

Configuration:

```text
RATE_LIMIT_ENABLED=true
RATE_LIMIT_STORE=redis
RATE_LIMIT_GLOBAL_RPM=300
RATE_LIMIT_AUTH_LOGIN_RPM=5
RATE_LIMIT_WRITE_RPM=60
```

## OWASP Top 10 Controls

ESSK must map implementation rules to OWASP Top 10 risks.

### A01 Broken Access Control

Required controls:

- Permission middleware for protected routes.
- Tenant boundary enforcement in service and repository queries.
- Deny by default when permission metadata is missing.
- Tests for cross-tenant access denial.

### A02 Cryptographic Failures

Required controls:

- Argon2id password hashing.
- Strong JWT signing secret or asymmetric keys.
- HTTPS in production.
- No secrets in logs or frontend bundles.
- Sensitive tokens stored hashed where persistence is required.

### A03 Injection

Required controls:

- sqlc parameterized queries only.
- No string-concatenated SQL for user input.
- Input validation for request body, params, and query string.
- Output encoding handled by frontend framework.

### A04 Insecure Design

Required controls:

- Threat-model checklist for auth, tenant, RBAC, and audit changes.
- Secure defaults in generated modules.
- Audit logging for security-relevant actions.
- Rate limiting on abuse-prone endpoints.

### A05 Security Misconfiguration

Required controls:

- Explicit CORS allowlist.
- Security headers.
- No debug mode in production.
- `.env.example` only; real secrets are not committed.
- Container runs as non-root.

### A06 Vulnerable And Outdated Components

Required controls:

- Dependency scanning in CI.
- `govulncheck` for Go.
- `pnpm audit` or equivalent for frontend.
- Dependabot or Renovate for dependency updates.

### A07 Identification And Authentication Failures

Required controls:

- Refresh token rotation.
- Short-lived access tokens.
- Account status checks.
- Login rate limiting.
- Password policy and secure reset flow when reset is added.

### A08 Software And Data Integrity Failures

Required controls:

- Lockfiles committed.
- CI uses frozen lockfile mode.
- Container image build provenance in future release phase.
- Migration review before production deployment.

### A09 Security Logging And Monitoring Failures

Required controls:

- Structured logs.
- Audit logs for auth, RBAC, tenant, and CRUD changes.
- Request ID and correlation ID.
- Security events visible in observability pipeline.

### A10 Server-Side Request Forgery

Required controls:

- No arbitrary URL fetch in foundation modules.
- Allowlist outbound integrations.
- Block private IP ranges for user-provided URLs when external fetch features are introduced.
- Timeout all outbound HTTP calls.

## Audit Logging

Audit these actions:

- Login success.
- Login failure.
- Logout.
- Refresh token rotation.
- User create/update/delete.
- Tenant create/update/delete.
- Role and permission changes.
- Access denied for critical permission failures.

Audit entry must include:

- Actor user ID.
- Tenant ID.
- Action.
- Resource type.
- Resource ID.
- IP address.
- User agent.
- Metadata.
- Timestamp.

## Observability Phase 1

Required:

- Structured JSON logs.
- Request ID propagated through logs and response headers.
- Health and readiness endpoints.
- Basic latency logging.

## Observability Future

Add OpenTelemetry when foundation is stable.

Future components:

- Traces: OpenTelemetry SDK.
- Metrics: Prometheus endpoint.
- Dashboards: Grafana.
- Logs: Loki or cloud provider logs.

Recommended metrics:

- HTTP request count.
- HTTP request duration.
- HTTP error count.
- Database query duration.
- Redis operation duration.
- Login success/failure count.

## Secret Management

Local:

- `.env` files.

CI:

- GitHub Actions secrets.

Production future:

- Kubernetes Secret.
- External Secrets Operator.
- Cloud secret manager.

Rules:

- No hardcoded secrets.
- No default production JWT secret.
- Rotate secrets through deployment process.
