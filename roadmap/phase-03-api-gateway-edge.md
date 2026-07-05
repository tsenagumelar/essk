# Phase 03 - API Gateway And Edge Layer

## Status

Current status: Not Started

## Goal

Introduce a production-style entry layer before backend services to support routing, security, rate limiting, and operational controls.

## Current State

- Backend is reachable directly on its exposed port.
- Docker Compose has local service wiring.
- Nginx has been considered in infrastructure specs, but a full gateway policy layer is not implemented.
- Request logging and rate limiting exist mostly inside backend code.

## Gap Analysis

| Area | Gap | Impact |
| --- | --- | --- |
| Gateway runtime | No dedicated gateway service is required for app flow | Production traffic path is not representative |
| Central routing | Routes are not centralized outside backend | Harder to introduce multiple services later |
| Gateway security | No gateway-level size limit, auth hook, API key hook, or CORS policy | Backend must handle every edge concern itself |
| Request metadata | No formal trusted headers contract | Correlation and tenant propagation can be inconsistent |
| CDN/WAF | Not implemented | No edge protection or static asset strategy |

## Tasks

| ID | Task | Status | Detail |
| --- | --- | --- | --- |
| 03-01 | Choose initial gateway | Not Started | Start with Nginx or Traefik; document ADR |
| 03-02 | Add gateway service to Compose | Not Started | Route web and API traffic through gateway |
| 03-03 | Add gateway route config | Not Started | Define `/api/v1/*`, web, health, and future internal route groups |
| 03-04 | Add request ID propagation | Not Started | Generate or forward `X-Request-ID` |
| 03-05 | Add gateway access logging | Not Started | Include request id, method, path, status, latency |
| 03-06 | Add gateway request limits | Not Started | Body size, timeout, header size |
| 03-07 | Add gateway CORS policy | Not Started | Centralize allowed origins and methods |
| 03-08 | Add gateway rate limit policy | Not Started | Protect auth and public APIs before backend |
| 03-09 | Define trusted header contract | Not Started | Tenant id, request id, forwarded host, protocol |
| 03-10 | Document CDN and WAF production strategy | Not Started | Cloudflare, AWS CloudFront/WAF, or equivalent |

## Exit Criteria

- Local production-like mode routes web and API through the gateway.
- Gateway config is version-controlled.
- Backend receives consistent request metadata.
- Gateway enforces basic traffic limits.
- Production CDN and WAF strategy is documented.

## Dependencies

- Phase 01 error and API consistency.
- Phase 02 tenant context contract.

## Notes

The gateway should be introduced before splitting services so future service extraction does not require changing client-facing routes.
