# ESSK Backend

Go backend for ESSK.

Current runtime mode:

- Separated service foundation under `services/*`.

The separated service foundation keeps one shared PostgreSQL database while reserving different schemas per service boundary. Service-to-service communication is prepared through gRPC and shared protobuf contracts.

## Run

```text
go run ./services/api-gateway
```

## Run Separated Services

Run from `services/backend`:

```text
go run ./services/auth-service
go run ./services/tenant-service
go run ./services/iam-service
go run ./services/catalog-service
go run ./services/audit-service
```

Default ports:

```text
api-gateway     HTTP 18080  gRPC 19100  schema gateway
auth-service    HTTP 18110  gRPC 19110  schema auth
tenant-service  HTTP 18120  gRPC 19120  schema tenant
iam-service     HTTP 18130  gRPC 19130  schema iam
catalog-service HTTP 18140  gRPC 19140  schema catalog
audit-service   HTTP 18150  gRPC 19150  schema audit
```

Common environment variables:

```text
ESSK_SERVICE_NAME=
ESSK_HTTP_PORT=
ESSK_GRPC_PORT=
ESSK_DB_SCHEMA=
ESSK_APP_ENV=
ESSK_APP_VERSION=
ESSK_LOG_LEVEL=
ESSK_LOG_PRETTY=
DATABASE_URL=
POSTGRES_HOST=
POSTGRES_PORT=
POSTGRES_DB=
POSTGRES_USER=
POSTGRES_PASSWORD=
```

Each separated service has this shape:

```text
services/{service-name}/
  main.go
  routes/
  handler/
  usecase/
  repositories/
```

Shared service runtime code lives in:

```text
internal/platform/
  config/
  helper/
  logger/
  service/
```

Shared protobuf and schema migrations live in:

```text
shared/
  protobuf/
  migrations/
```

## Migrations

Run from `services/backend`:

```text
go run ./cmd/essk migrate up
```

Run shared service schema migrations:

```text
go run ./cmd/essk migrate shared up
```

Shared schema migrations use a separate migration table:

```text
shared_schema_migrations
```

## Seed Admin

Run after migrations:

```text
go run ./cmd/essk seed admin
```

Default local credentials:

```text
admin@essk.local
Admin123!
```

## Policy Check

```text
go run ./cmd/essk policy check --all
```
