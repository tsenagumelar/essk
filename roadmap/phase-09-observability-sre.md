# Phase 09 - Observability And SRE Readiness

## Status

Current status: Not Started

## Goal

Make the platform observable, measurable, and operable.

## Current State

- Backend uses structured logging.
- Health/readiness foundations exist.
- No full tracing, metrics, dashboards, alerts, or runbooks exist.

## Gap Analysis

| Area | Gap | Impact |
| --- | --- | --- |
| Tracing | OpenTelemetry missing | Cross-service and worker debugging is hard |
| Metrics | Prometheus metrics missing | No objective service health view |
| Dashboards | Grafana dashboards missing | Operators lack real-time visibility |
| Alerts | Alert rules missing | Incidents may be discovered late |
| Runbooks | Missing | Operational response is inconsistent |
| SLOs | Partially documented, not enforced | Reliability goals are not measurable |

## Tasks

| ID | Task | Status | Detail |
| --- | --- | --- | --- |
| 09-01 | Add OpenTelemetry tracing | Not Started | API, DB, Redis, event bus, workers |
| 09-02 | Add Prometheus metrics endpoint | Not Started | Request count, latency, errors, DB stats, worker stats |
| 09-03 | Add Grafana to Compose profile | Not Started | Local observability stack |
| 09-04 | Add default dashboards | Not Started | API, database, Redis, event bus, workers |
| 09-05 | Add correlation id across layers | Not Started | Gateway, API, events, workers, logs |
| 09-06 | Add alert rules | Not Started | Error rate, latency, database down, queue lag, worker failures |
| 09-07 | Define SLOs per service | Not Started | Availability, latency, error rate |
| 09-08 | Add error budget policy | Not Started | Define release behavior when reliability is poor |
| 09-09 | Add incident runbooks | Not Started | API degraded, DB degraded, event backlog, worker failure, restore |
| 09-10 | Add audit/security log dashboard | Not Started | Track failed login, permission changes, suspicious API key usage |

## Exit Criteria

- Operators can inspect latency, errors, throughput, DB health, queue lag, and worker failures.
- Logs, traces, and metrics share correlation ids.
- Common incidents have documented runbooks.
- Service reliability can be measured against SLOs.

## Dependencies

- Phase 03 gateway for edge metrics.
- Phase 04 event bus.
- Phase 05 workers.
- Phase 08 database scale.

## Notes

Observability should be added before Kubernetes production rollout. Otherwise debugging HA deployment issues will be unnecessarily difficult.
