# ESSK Enterprise Platform Roadmap

## Purpose

This folder breaks down the enterprise platform architecture roadmap into actionable phase documents.

Main architecture reference:

```text
specs/15-enterprise-platform-architecture-roadmap.md
specs/architecture.png
```

The roadmap follows the provided Enterprise Platform Architecture diagram as the long-term target.

## Status Legend

| Status | Meaning |
| --- | --- |
| Not Started | No implementation work has been done yet |
| In Progress | Some implementation exists, but the phase is incomplete |
| Partially Done | A usable foundation exists, but enterprise-level requirements are still missing |
| Done | Requirements and exit criteria are complete |
| Blocked | Work cannot continue until a dependency is resolved |

## Phase Index

| Phase | Document | Current Status | Goal |
| --- | --- | --- | --- |
| 01 | [Foundation Stabilization](./phase-01-foundation-stabilization.md) | In Progress | Stabilize modular backend/frontend, contracts, tests, and standards |
| 02 | [Enterprise Multi-Tenancy And Security](./phase-02-enterprise-multitenancy-security.md) | Partially Done | Enforce tenant isolation and enterprise security |
| 03 | [API Gateway And Edge Layer](./phase-03-api-gateway-edge.md) | Not Started | Add production-style gateway and edge controls |
| 04 | [Event-Driven Foundation](./phase-04-event-driven-foundation.md) | Not Started | Add outbox, event bus, and event contracts |
| 05 | [Worker Services](./phase-05-worker-services.md) | Not Started | Move async workloads to independently scalable workers |
| 06 | [Platform Services](./phase-06-platform-services.md) | Not Started | Add notification, file, search, workflow, audit, and plugin services |
| 07 | [Module Scaffolding And Plugin Model](./phase-07-module-scaffolding-plugin-model.md) | Not Started | Generate complete modules consistently |
| 08 | [Database Scale Strategy](./phase-08-database-scale-strategy.md) | Not Started | Prepare database for large tenants and large tables |
| 09 | [Observability And SRE Readiness](./phase-09-observability-sre.md) | Not Started | Add tracing, metrics, dashboards, alerts, and runbooks |
| 10 | [Kubernetes And High Availability](./phase-10-kubernetes-high-availability.md) | Not Started | Support HA production deployment |
| 11 | [Enterprise Integration And Ecosystem](./phase-11-enterprise-integration-ecosystem.md) | Not Started | Add external integration, SDK, webhook, and plugin ecosystem |

## Recommended Execution Order

1. Complete Phase 01.
2. Complete Phase 02 security-critical work.
3. Add Phase 03 gateway foundation.
4. Add Phase 04 outbox and event bus.
5. Add Phase 05 first workers.
6. Continue Phase 07 scaffolding once module structure is stable.
7. Add Phase 06 platform services incrementally.
8. Add Phase 08 database scale improvements.
9. Add Phase 09 observability before production readiness.
10. Add Phase 10 Kubernetes deployment.
11. Add Phase 11 ecosystem features.

## Governance

- Every phase document must keep `Status`, `Gap Analysis`, `Tasks`, and `Exit Criteria` updated.
- A task can only move to `Done` when implementation, tests, and documentation are complete.
- Architecture changes that affect multiple phases must update `specs/15-enterprise-platform-architecture-roadmap.md`.
- Major decisions must be documented as ADRs.
