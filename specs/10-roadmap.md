# ESSK Specs - Roadmap

## Milestone 1 - Foundation

Goal: one-command local startup and basic full-stack connectivity.

Backend:

- Go service bootstrap.
- Fiber router.
- Config loader.
- Logger.
- PostgreSQL connection.
- Redis connection.
- Health and readiness endpoints.
- Standard response helper.
- Global error handler.
- Request validator.
- Swagger.
- Migration command.

Frontend:

- Next.js app bootstrap.
- Tailwind setup.
- Login page UI.
- Authenticated layout shell.
- Dashboard placeholder.
- API health check screen.

Infrastructure:

- Docker Compose.
- Backend Dockerfile.
- Web Dockerfile.
- PostgreSQL service.
- Redis service.
- Nginx reverse proxy.
- `.env.example` files.

CI:

- Backend test/build.
- Frontend lint/typecheck/build.
- Docker build validation.

Exit criteria:

- `docker compose up` starts all foundation services.
- `GET /health` works.
- Frontend can call backend.
- Swagger is available.
- Migration command works.

## Milestone 2 - Auth Core

Backend:

- User table.
- Refresh token table.
- Password hashing.
- Login.
- Logout.
- Refresh token rotation.
- Auth middleware.
- `/auth/me`.
- Seed admin user.

Frontend:

- Login integration.
- Token handling.
- Protected routes.
- Logout.
- Profile display.

Exit criteria:

- Developer can login with seeded admin.
- Protected API rejects unauthenticated requests.
- Protected UI redirects unauthenticated users.

## Milestone 3 - Tenant And RBAC

Backend:

- Tenant CRUD.
- Role CRUD.
- Permission list.
- Role-permission assignment.
- User-role assignment.
- Permission middleware.
- Tenant boundary enforcement.

Frontend/Admin:

- Tenant list and form.
- User list and form.
- Role list and permission assignment.

Exit criteria:

- Admin can manage users, tenants, roles, and permissions.
- Tenant-scoped users cannot access other tenant data.

## Milestone 4 - Audit And Hardening

Backend:

- Audit log table.
- Audit writer.
- Audit endpoints.
- Rate limiting for auth.
- Stronger security headers.
- Integration tests for tenant and RBAC behavior.

Frontend/Admin:

- Audit log viewer.
- Filtering by actor, action, resource, date.

Exit criteria:

- Critical actions produce audit records.
- Admin can inspect audit logs.
- CI runs meaningful test suite.

## Milestone 5 - Product Module Template

Goal: make starter kit easy to extend into real SaaS products.

Deliverables:

- `essk add-module {module-name}` command.
- Backend module generator.
- Frontend CRUD page generator.
- Migration and query generator.
- API documentation example.
- Test template.
- Example generated module.

Exit criteria:

- Developer can run `essk add-module user-management` and get backend CRUD endpoint files, frontend CRUD page files, migration files, query files, and starter tests.
- Generated files compile after codegen and formatting.
- Developer can create a new module consistently in less than one hour.

## Future Platform Modules

- Workflow Engine.
- Approval Engine.
- Notification Center.
- Scheduler.
- File Storage.
- Report Engine.
- API Gateway.
- Feature Flags.
- Configuration Center.
- Secret Manager.
- Observability stack.
- AI Agent integration.
