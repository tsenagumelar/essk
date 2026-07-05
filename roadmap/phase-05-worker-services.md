# Phase 05 - Worker Services

## Status

Current status: Not Started

## Goal

Move slow, heavy, scheduled, and cross-cutting workloads out of synchronous API requests.

## Current State

- No worker runtime exists.
- Backend runs as a single API service.
- No separate worker deployment or command exists.
- No queue consumer implementation exists.

## Gap Analysis

| Area | Gap | Impact |
| --- | --- | --- |
| Worker runtime | Missing | Async workloads cannot scale independently |
| Worker health | Missing | Operators cannot monitor workers separately |
| Audit worker | Missing | Audit processing cannot be decoupled |
| Notification worker | Missing | Email/SMS/push cannot be retried safely |
| Report worker | Missing | Exports can block API requests if implemented synchronously |
| Scheduler | Missing | Recurring jobs have no platform home |

## Tasks

| ID | Task | Status | Detail |
| --- | --- | --- | --- |
| 05-01 | Add worker command | Not Started | Example: `go run ./cmd/essk worker` |
| 05-02 | Add worker configuration | Not Started | Concurrency, queues, retry, shutdown timeout |
| 05-03 | Add graceful shutdown | Not Started | Stop accepting jobs and finish in-flight work |
| 05-04 | Add worker heartbeat | Not Started | Expose health or write heartbeat status |
| 05-05 | Add audit worker | Not Started | Consume audit events and write audit logs |
| 05-06 | Add notification worker | Not Started | Consume notification events and call provider adapter |
| 05-07 | Add report worker | Not Started | Generate export files asynchronously |
| 05-08 | Add schedule worker | Not Started | Run recurring platform jobs |
| 05-09 | Add integration worker | Not Started | Handle webhooks and third-party sync jobs |
| 05-10 | Add workflow worker skeleton | Not Started | Progress workflow steps from events |
| 05-11 | Add worker observability fields | Not Started | Job id, event id, tenant id, duration, status, error |

## Exit Criteria

- API and workers can run as separate processes.
- Workers can be scaled independently.
- Failed jobs can retry and be inspected.
- At least audit and notification worker flows exist.

## Dependencies

- Phase 04 event-driven foundation.
- Phase 09 observability will improve this phase later.

## Notes

Workers should consume from the same event envelope standard defined in Phase 04.
