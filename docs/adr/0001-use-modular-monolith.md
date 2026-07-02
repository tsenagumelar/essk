# ADR 0001: Use Modular Monolith

## Status

Accepted

## Context

ESSK is a starter kit for multiple enterprise SaaS products. The foundation needs strong module boundaries without the operational overhead of distributed services during the early product phase.

## Decision

Use a modular monolith for the backend foundation. Each business capability is implemented as an internal module with consistent handler, command, query, service, repository, validation, and test boundaries.

## Consequences

- Local development stays simple.
- Cross-module refactoring remains practical.
- Module boundaries can later guide service extraction.
- Teams must enforce boundaries through code review and policy checks.

## Alternatives Considered

- Microservices from the start.
- Single unstructured monolith.
