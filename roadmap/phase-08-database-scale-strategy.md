# Phase 08 - Database Scale Strategy

## Status

Current status: Not Started

## Goal

Prepare PostgreSQL for large tenants, large tables, reporting workloads, and enterprise data governance.

## Current State

- PostgreSQL runs locally through Docker Compose.
- Tables use audit fields and soft delete conventions.
- Tenant-related tables exist.
- No read replica, PgBouncer, partitioning, RLS, or archive strategy exists yet.

## Gap Analysis

| Area | Gap | Impact |
| --- | --- | --- |
| Connection pooling | PgBouncer missing | High connection count can overload PostgreSQL |
| Read replicas | Missing | Reporting and heavy reads compete with writes |
| Partitioning | Missing | Large audit/log tables can degrade over time |
| Archiving | Missing | Old data can make operational tables too large |
| RLS | Missing | Tenant isolation relies on application code |
| Query policy | Needs formal review process | Slow queries can reach production |

## Tasks

| ID | Task | Status | Detail |
| --- | --- | --- | --- |
| 08-01 | Add PgBouncer local option | Not Started | Compose profile for connection pooling |
| 08-02 | Define index review checklist | Not Started | Required indexes for tenant, status, created date, foreign keys |
| 08-03 | Define required tenant table policy | Not Started | Tenant-scoped tables must include `tenant_id` unless ADR says otherwise |
| 08-04 | Add RLS policy template | Not Started | Standard PostgreSQL tenant isolation policy |
| 08-05 | Apply RLS to tenant-scoped tables | Not Started | Users, roles, products, future generated modules |
| 08-06 | Define partitioning strategy | Not Started | Audit logs, notifications, integration logs, activity logs |
| 08-07 | Add archive strategy | Not Started | Move old data to archive tables or object storage |
| 08-08 | Define read replica strategy | Not Started | Reporting and analytics use read replicas |
| 08-09 | Add read/write repository pattern | Not Started | Query service can read from replica where safe |
| 08-10 | Add slow query logging guide | Not Started | Enable and monitor slow SQL in production |
| 08-11 | Add database migration review checklist | Not Started | Prevent unsafe migrations |
| 08-12 | Add query performance tests | Not Started | Test critical list/search endpoints with larger data |

## Exit Criteria

- Tenant-scoped data has database-level isolation.
- Large append-only tables have partition and retention plans.
- Read-heavy workloads have a read replica path.
- Connection pooling is available for production-like deployment.
- Database migrations follow a safety checklist.

## Dependencies

- Phase 02 tenant and security model.
- Phase 04/05 for event and worker tables.
- Phase 06 for notification, file, search, workflow tables.

## Notes

This phase is required before claiming the platform can support very large tenants or millions of rows per module.
