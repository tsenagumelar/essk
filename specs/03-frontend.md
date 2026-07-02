# ESSK Specs - Frontend

## Stack

- Next.js 15+ App Router.
- React 19+.
- TypeScript.
- Tailwind CSS.
- Radix UI primitives.
- shadcn/ui style component organization.
- TanStack Query for server state.
- Zustand for small client-only state.
- React Hook Form for forms.
- Zod for schema validation.
- axios or generated fetch client for HTTP.
- Vitest for unit tests.
- React Testing Library for component tests.
- Playwright for end-to-end tests.

## Web App Structure

```text
apps/web/
  app/
    (auth)/
      login/
    (app)/
      dashboard/
      profile/
    layout.tsx
    page.tsx

  components/
    ui/
    layout/
    forms/
    data-table/

  features/
    auth/
    user/
    tenant/
    rbac/

  lib/
    api/
    auth/
    query/
    env/
    utils/

  stores/
  styles/
  tests/
```

## App Responsibilities

Foundation web app must provide:

- Login page.
- Logout action.
- Authenticated layout.
- Dashboard placeholder.
- Profile page.
- API health display.
- Protected route handling.
- Global loading and error states.

Admin app later provides:

- User CRUD.
- Tenant CRUD.
- Role CRUD.
- Permission management.
- Audit log list.

## Rendering Strategy

Use server components by default.

Use client components for:

- Forms.
- Mutations.
- Interactive tables.
- Menus/dialogs.
- Client-side auth state.

Authentication-dependent pages should validate session/token before rendering protected content.

## API Integration

Preferred flow:

1. Backend exposes OpenAPI spec.
2. `packages/sdk-ts` generates TypeScript types.
3. `apps/web` consumes SDK through `lib/api`.
4. TanStack Query wraps read and mutation operations.

API client requirements:

- Base URL from environment.
- Request ID support if manually provided.
- Access token attached to protected requests.
- Refresh flow handled centrally.
- Standard response unwrapped consistently.
- Errors normalized into typed frontend errors.

## State Management

TanStack Query:

- Server data.
- Cache invalidation.
- Pagination.
- Mutations.

Zustand:

- Sidebar state.
- Theme preference.
- Lightweight UI state.

Do not duplicate server data in Zustand.

## Forms

Use React Hook Form + Zod.

Each form feature should include:

- Zod schema.
- Type inferred from schema.
- Field-level error rendering.
- Submit loading state.
- Backend validation error mapping.

## Frontend Module Scaffolding

Frontend must support the same module scaffolding command used by backend.

Command:

```text
essk add-module user-management
```

Frontend output for `user-management`:

```text
apps/web/app/(app)/user-management/
  page.tsx
  loading.tsx
  error.tsx
  new/
    page.tsx
  [id]/
    page.tsx
    edit/
      page.tsx

apps/web/features/user-management/
  api.ts
  schema.ts
  types.ts
  hooks.ts
  components/
    user-management-form.tsx
    user-management-table.tsx
    user-management-delete-dialog.tsx
```

Generated frontend CRUD pages:

- List page with table, search, pagination, create button, edit action, delete action.
- Create page with generated form.
- Detail page.
- Edit page.
- Delete confirmation dialog.

Generated frontend integrations:

- Zod schema.
- React Hook Form setup.
- TanStack Query hooks for list, detail, create, update, and delete.
- API wrapper using `packages/sdk-ts` when available.
- Fallback API wrapper using the standard HTTP client.

Scaffolding options:

```text
essk add-module user-management --tenant-scoped=true --with-frontend=true --with-backend=true
```

Rules:

- Frontend module route path uses kebab-case.
- Feature folder uses kebab-case.
- Component filenames use kebab-case.
- Generated pages must use the authenticated app layout.
- Generated code must include loading and error states.
- Generated CRUD must use standard API response handling.

## Styling

Use Tailwind CSS tokens through CSS variables:

- Background.
- Foreground.
- Muted.
- Border.
- Primary.
- Secondary.
- Destructive.
- Ring.

Design preference:

- Dense enterprise UI.
- Clear hierarchy.
- Accessible contrast.
- Minimal decorative elements.
- Stable responsive layouts.

## Environment Variables

Required:

```text
NEXT_PUBLIC_APP_NAME=ESSK
NEXT_PUBLIC_API_BASE_URL=http://localhost:18080/api/v1
```

Do not expose secrets through `NEXT_PUBLIC_*`.

## Testing

Unit/component:

- Form validation.
- Auth state behavior.
- API error mapping.
- Component rendering states.

E2E:

- Login.
- Protected page access.
- Logout.
- API health page.

## Quality Gates

Required commands:

```text
pnpm lint
pnpm typecheck
pnpm test
pnpm e2e
```
