# ESSK Specs - Enterprise Readiness

## Purpose

This spec captures enterprise-level requirements that sit above individual backend, frontend, database, and infrastructure implementation details. These rules prevent the starter kit from becoming only a working demo; it must remain suitable as the foundation for regulated, multi-tenant, long-lived SaaS products.

## Architecture Decision Records

Every significant technical decision must be documented as an ADR.

Location:

```text
docs/adr/
```

Naming:

```text
0001-use-sqlc-and-pgx.md
0002-use-modular-monolith.md
0003-use-redis-backed-rate-limiter.md
```

ADR template:

```text
# ADR {number}: {title}

## Status

Accepted | Proposed | Deprecated | Superseded

## Context

## Decision

## Consequences

## Alternatives Considered
```

Required ADR topics:

- Backend framework.
- Query layer.
- Authentication model.
- Multi-tenancy model.
- Rate limiting strategy.
- CQRS approach.
- Deployment model.
- Module scaffolding strategy.

## API Lifecycle

Rules:

- Public API is versioned under `/api/v1`.
- Breaking changes require a new version or explicit migration plan.
- Deprecated endpoints must return deprecation metadata before removal.
- OpenAPI must be updated in the same change as API behavior.
- Generated SDKs must be regenerated after OpenAPI changes.

Required headers for deprecated endpoints:

```text
Deprecation: true
Sunset: {date}
Link: <{migration-doc-url}>; rel="deprecation"
```

## Service Level Objectives

Default foundation SLO targets:

- API availability target: 99.9% for production deployments.
- API p95 latency target: under 300 ms for common CRUD reads.
- API p95 latency target: under 500 ms for common CRUD writes.
- Login p95 latency target: under 700 ms.
- Frontend initial route load target: under 2.5 seconds on normal broadband.

These are starter targets. Product teams may tighten or relax them through documented ADRs.

## Performance Budgets

Backend:

- Pagination is mandatory for list endpoints.
- Default page size is 20.
- Maximum page size is 100.
- Slow query logging must be enabled in production.
- Queries expected to return large datasets must use indexes and projections.

Frontend:

- Avoid large client-only bundles for enterprise CRUD screens.
- Use route-level code splitting.
- Avoid duplicating server state in client stores.
- Data tables must support pagination.

## Backup And Restore

PostgreSQL backup requirements:

- Daily logical backup in production.
- Point-in-time recovery target when infrastructure supports it.
- Restore procedure documented.
- Restore test at least once per release cycle for production-grade environments.

Redis backup requirements:

- Redis is treated as disposable cache in the foundation.
- Persistent Redis usage must be documented before production use.

Required documentation:

```text
docs/operations/backup-restore.md
```

## Disaster Recovery

Minimum DR documentation:

- Recovery Time Objective.
- Recovery Point Objective.
- Database restore steps.
- Secret rotation steps.
- Rollback steps.
- Communication checklist.

Default starter targets:

- RTO: 4 hours.
- RPO: 24 hours.

Product teams must revise these targets before production launch.

## Data Retention

Default retention rules:

- Audit logs: retain for at least 1 year.
- Refresh tokens: delete expired and revoked records after 30 days.
- Soft-deleted business records: retained until product-specific policy is defined.
- Application logs: retain according to environment and compliance requirements.

Every module handling business data must document:

- Retention period.
- Deletion behavior.
- Export behavior.
- Audit requirements.

## Privacy And Compliance

Foundation requirements:

- Avoid storing unnecessary personal data.
- Mark personally identifiable data fields in module docs.
- Do not log personal data unless explicitly required.
- Support user data export at product level when required.
- Support deletion/anonymization workflows at product level when required.

Compliance readiness topics:

- SOC 2 readiness.
- ISO 27001 readiness.
- GDPR-style data subject request support.
- Auditability for admin actions.

The starter kit does not claim certification by itself; it provides patterns that make certification work easier.

## Dependency Governance

Rules:

- Lockfiles are committed.
- CI uses frozen lockfile mode.
- Dependencies must be scanned for vulnerabilities.
- Licenses must be checked before release.
- Direct dependencies should be reviewed before adoption.
- Avoid adding libraries for trivial helpers.

Recommended tools:

- `govulncheck`.
- `pnpm audit`.
- Dependabot or Renovate.
- SBOM generation in release workflow.

Policy command:

```text
essk dependency check
essk license check
```

## Configuration Governance

Rules:

- Every environment variable must be documented in `.env.example`.
- Production must not rely on insecure defaults.
- Feature flags must default to the safest behavior.
- Secrets must never use `NEXT_PUBLIC_*`.
- Config parsing must fail fast for missing required production values.

Required docs:

```text
docs/configuration.md
```

## Feature Flags

Feature flags are optional in foundation but the design must reserve a path for them.

Recommended future approach:

- Database-backed feature flags for product behavior.
- Environment-backed flags for infrastructure behavior.
- Tenant-scoped flags where product modules require gradual rollout.

Rules:

- Feature flags must not bypass authorization.
- Feature flags must be auditable when changed.
- Stale flags must be removed.

## Release And Rollback

Every release must define:

- Version.
- Migration list.
- Feature changes.
- API changes.
- Operational notes.
- Rollback path.

Rollback rules:

- Application rollback must be possible independently from database rollback where feasible.
- Destructive migrations require explicit approval.
- Backward-compatible migrations are preferred.
- Release notes must include smoke test steps.

Required docs:

```text
docs/operations/release.md
docs/operations/rollback.md
```

## Operational Runbooks

Required runbooks before production use:

- Local development.
- Deployment.
- Rollback.
- Backup and restore.
- Incident response.
- Secret rotation.
- Database migration.
- Creating a new module.

Location:

```text
docs/operations/
```

## Incident Response

Minimum incident workflow:

1. Detect.
2. Triage.
3. Mitigate.
4. Communicate.
5. Resolve.
6. Review.

Required post-incident output:

- Timeline.
- Root cause.
- Impact.
- Corrective actions.
- Preventive actions.

## Enterprise Readiness Checklist

Before using ESSK as a production product foundation, the repository must have:

- ADRs for major decisions.
- OWASP Top 10 controls mapped.
- Rate limiter enabled.
- CQRS module template implemented.
- Pre-commit and CI policy checks.
- Backup and restore documentation.
- Release and rollback documentation.
- Dependency and license checks.
- OpenAPI generation and SDK generation.
- Operational runbooks.
