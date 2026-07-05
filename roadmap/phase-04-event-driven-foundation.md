# Phase 04 - Event-Driven Foundation

## Status

Current status: Not Started

## Goal

Add asynchronous event publishing and processing foundations without prematurely splitting core services into microservices.

## Current State

- No event bus exists.
- No outbox table exists.
- No event envelope standard exists.
- Audit currently runs synchronously or within application logic depending on module behavior.

## Gap Analysis

| Area | Gap | Impact |
| --- | --- | --- |
| Event bus | Kafka/NATS not implemented | No async module communication |
| Outbox | No transactional outbox | Risk of data saved without event or event sent without data |
| Event schema | No versioned event contract | Future consumers can break silently |
| Retry/DLQ | No retry or dead letter pattern | Failed async work can be lost |
| Idempotency | No event idempotency standard | Consumers may duplicate side effects |

## Tasks

| ID | Task | Status | Detail |
| --- | --- | --- | --- |
| 04-01 | Choose first event bus | Not Started | Prefer NATS first; document Kafka upgrade path |
| 04-02 | Add event bus to Compose | Not Started | Local developer environment must include event bus |
| 04-03 | Create outbox migration | Not Started | Store event envelope in PostgreSQL |
| 04-04 | Add outbox repository | Not Started | Write and claim events safely |
| 04-05 | Add outbox publisher command | Not Started | Publish pending events to bus |
| 04-06 | Define event envelope | Not Started | Include id, type, version, tenant, actor, correlation, causation, payload, occurred at |
| 04-07 | Define topic naming convention | Not Started | Example: `essk.user.created.v1` |
| 04-08 | Add retry and backoff strategy | Not Started | Track attempt count and next retry time |
| 04-09 | Add dead letter queue strategy | Not Started | Store or publish failed events after retry exhaustion |
| 04-10 | Add idempotency helper | Not Started | Consumers can record processed event ids |
| 04-11 | Publish first user and tenant events | Not Started | `user.created`, `user.updated`, `tenant.created`, `role.assigned` |

## Exit Criteria

- Business transactions can write data and outbox events atomically.
- Events are published asynchronously.
- Event payloads are versioned.
- Failed publishing can retry and eventually move to DLQ.
- Consumers can be idempotent.

## Dependencies

- Phase 01 backend module consistency.
- Phase 02 tenant context and audit policy.

## Notes

This phase is the bridge from modular monolith toward event-driven platform architecture.
