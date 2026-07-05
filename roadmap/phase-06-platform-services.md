# Phase 06 - Platform Services

## Status

Current status: Not Started

## Goal

Implement shared enterprise platform services that product modules can reuse.

## Current State

- Auth, tenant, RBAC, user, product, and audit foundations exist.
- No notification, file, search, workflow, or plugin service exists.
- Audit exists but needs stronger platform-service treatment.

## Gap Analysis

| Area | Gap | Impact |
| --- | --- | --- |
| Notification | No provider adapters or templates | Product modules cannot send reliable messages |
| File storage | No S3-compatible storage | Upload/download features cannot be standardized |
| Search | No search engine | Large text search and analytics are limited |
| Workflow | No approval engine | Enterprise approval processes cannot be modeled |
| Plugin registry | No plugin service | Modules cannot be enabled/disabled per tenant |
| Audit service | Needs stronger query/export/retention support | Compliance use cases remain incomplete |

## Tasks

| ID | Task | Status | Detail |
| --- | --- | --- | --- |
| 06-01 | Add notification service model | Not Started | Notification, template, channel, delivery attempt |
| 06-02 | Add email provider adapter | Not Started | SMTP or provider abstraction |
| 06-03 | Add SMS/WhatsApp adapter placeholder | Not Started | Provider-neutral interface |
| 06-04 | Add file service model | Not Started | File metadata, owner, tenant, storage key, checksum |
| 06-05 | Add MinIO to Compose | Not Started | Local S3-compatible storage |
| 06-06 | Add file upload/download API | Not Started | Tenant-aware file operations |
| 06-07 | Add search engine to Compose | Not Started | OpenSearch or Elasticsearch |
| 06-08 | Add indexing worker | Not Started | Index tenant-aware resources asynchronously |
| 06-09 | Add search API | Not Started | Shared query endpoint and module-specific adapters |
| 06-10 | Add workflow domain model | Not Started | Workflow, step, approver, instance, action |
| 06-11 | Add workflow API skeleton | Not Started | Start, approve, reject, cancel |
| 06-12 | Add plugin registry model | Not Started | Plugin manifest, module code, tenant enablement |
| 06-13 | Expand audit query/export | Not Started | Filtering, export, retention policy |

## Exit Criteria

- Product modules can call notification, file, search, workflow, audit, and plugin services through stable internal APIs.
- Platform services are tenant-aware and audit-aware.
- Local Compose can run the required development services.

## Dependencies

- Phase 04 event bus.
- Phase 05 workers.
- Phase 08 database scale for large audit/search/file metadata tables.

## Notes

These services should remain inside the modular platform first. Extract only when operational or scaling needs justify it.
