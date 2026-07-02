# ADR 0002: Use sqlc And pgx

## Status

Accepted

## Context

Enterprise SaaS products usually need explicit query control, predictable performance, and strong typing for database access.

## Decision

Use SQL written in query files, generate Go types with sqlc, and execute queries through pgx.

## Consequences

- SQL remains visible and reviewable.
- Generated types reduce runtime mapping errors.
- Query optimization stays straightforward.
- Developers must maintain SQL files and generated code.

## Alternatives Considered

- GORM.
- Hand-written database scanning only.
