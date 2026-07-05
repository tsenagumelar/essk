# ESSK Specs - Web Structure

## Purpose

The web app must be easy to maintain, readable, scalable, and safe to extend.

The target structure separates routing, feature implementation, shared UI, shared logic, and low-level infrastructure. New features should be added by creating a new feature folder and a route entry only, without modifying unrelated modules.

## Core Rules

- `apps/web/app` is only for routing, route metadata, layouts, loading states, and error boundaries.
- Route UI and business logic live in `apps/web/features/{feature}`.
- Reusable UI, hooks, utilities, and contracts live in `apps/web/shared`.
- Shared `.tsx` views must live under `apps/web/shared/components`.
- Shared hooks must live under `apps/web/shared/hooks` and should avoid rendering JSX when a shared component can render it instead.
- Shared non-view functions, API clients, session helpers, formatters, and export helpers must live under `apps/web/shared/functions`.
- Atomic design is used for reusable visual components.
- Each feature owns its view, types, hooks, API adapters, schemas, and feature-specific components.
- Cross-feature imports are not allowed unless routed through `shared` or an explicit public feature API.
- Feature folders expose a small public surface through `index.tsx` or `index.ts`.
- Low-level infrastructure stays in `apps/web/lib`.

## Target Directory Structure

```text
apps/web/
  app/
    layout.tsx
    page.tsx
    globals.css

    (auth)/
      login/
        page.tsx
        loading.tsx
        error.tsx

    (app)/
      layout.tsx
      dashboard/
        page.tsx
      health/
        page.tsx
      products/
        page.tsx
      profile/
        page.tsx
      tenants/
        page.tsx
      users/
        page.tsx
      roles/
        page.tsx

  features/
    auth/
      index.tsx
      types.ts
      hooks.ts
      api.ts
      schema.ts
      components/

    app-shell/
      index.tsx
      types.ts
      hooks.ts
      components/

    products/
      index.tsx
      types.ts
      hooks.ts
      api.ts
      schema.ts
      components/

    tenants/
      index.tsx
      types.ts
      hooks.ts
      api.ts
      schema.ts
      components/

    users/
      index.tsx
      types.ts
      hooks.ts
      api.ts
      schema.ts
      components/

    roles/
      index.tsx
      types.ts
      hooks.ts
      api.ts
      schema.ts
      components/

    system/
      index.tsx
      types.ts
      hooks.ts
      api.ts
      components/

  shared/
    components/
      atoms/
        button.tsx
        icon-button.tsx
        link-button.tsx
        text-input.tsx
        textarea-field.tsx
        select-field.tsx
        checkbox-field.tsx
        radio-group.tsx
        switch-field.tsx
        label.tsx
        helper-text.tsx
        badge.tsx
        avatar.tsx
        spinner.tsx
        skeleton.tsx
        separator.tsx
        kbd.tsx

      molecules/
        search-box.tsx
        filter-select.tsx
        multi-select.tsx
        pagination.tsx
        confirmation-dialog.tsx
        confirmable-action-dialog.tsx
        modal.tsx
        row-actions.tsx
        alert.tsx
        empty-state.tsx
        loading-state.tsx
        description-list.tsx
        dropdown-menu.tsx
        tabs.tsx
        breadcrumbs.tsx
        status-badge.tsx

      organisms/
        data-table.tsx
        crud-toolbar.tsx
        page-shell.tsx
        section-header.tsx
        stat-grid.tsx
        form-panel.tsx
        detail-panel.tsx
        app-navbar.tsx
        app-sidebar.tsx

      templates/
        crud-page-template.tsx
        authenticated-layout.tsx
        split-auth-layout.tsx

    hooks/
      use-confirmable-action.ts
      use-debounced-value.ts
      use-pagination.ts
      use-unauthorized-redirect.ts

    functions/
      api/
        client.ts
        errors.ts
        envelope.ts

      auth/
        session.ts
        permissions.ts

      export/
        export-excel.ts
        print-pdf.ts

      format/
        currency.ts
        date.ts

    types/
      api.ts
      pagination.ts
      table.ts

  lib/
    env.ts
    query/
      query-provider.tsx

  public/
```

## App Folder Contract

The `app` folder must stay thin.

Allowed in route files:

- Import a feature public view.
- Pass route params/search params to the feature.
- Export route metadata when needed.
- Define `loading.tsx` and `error.tsx`.

Not allowed in route files:

- API calls.
- Mutation logic.
- Table logic.
- Form state.
- Business rules.
- Feature-specific JSX beyond route composition.

Example:

```tsx
import { UsersView } from "@/features/users";

export default function UsersPage() {
  return <UsersView />;
}
```

## Feature Module Contract

Each route-backed feature should follow this structure:

```text
features/{feature}/
  index.tsx
  types.ts
  hooks.ts
  api.ts
  schema.ts
  components/
    {feature}-form.tsx
    {feature}-table.tsx
    {feature}-filters.tsx
```

### `index.tsx`

Owns the feature view composition.

Responsibilities:

- Compose feature hooks and components.
- Render loading, error, empty, and success states.
- Use shared organisms such as `DataToolbar`, `DataTable`, and `ModalForm`.
- Export the route-facing view, for example `UsersView`.

### `types.ts`

Owns feature-specific TypeScript types.

Examples:

- Entity response types.
- Form state types.
- Filter state types.
- Table row view models.

### `hooks.ts`

Owns feature-specific data and UI logic.

Examples:

- `useUsers`
- `useUserMutations`
- `useUserFilters`
- `useUserTable`
- `useUserForm`

Rules:

- Hooks can import feature `api.ts`.
- Hooks can import shared hooks.
- Hooks should keep `index.tsx` readable.

### `api.ts`

Owns feature API calls.

Rules:

- Uses shared API client only.
- Does not know about UI.
- Returns typed data.
- Does not directly touch local storage or router.

### `schema.ts`

Owns validation schemas.

Preferred libraries:

- Zod for runtime schema and form validation.
- React Hook Form for form integration.

### `components/`

Owns feature-specific components only.

Examples:

- `users-form.tsx`
- `users-table.tsx`
- `users-filters.tsx`
- `role-badges.tsx`

If a component becomes reusable by two or more features, move it to `shared`.

## Atomic Design Contract

Atomic design is applied only to reusable shared UI.

### Atoms

Small, style-consistent primitives.

Examples:

- `Button`
- `Input`
- `Select`
- `Checkbox`
- `Badge`
- `IconButton`

Rules:

- No business logic.
- No API calls.
- No feature imports.
- Minimal state only when required for accessibility.

### Molecules

Small reusable combinations of atoms.

Examples:

- `SearchBox`
- `FilterSelect`
- `Pagination`
- `ConfirmationDialog`
- `DropdownCheckbox`

Rules:

- No feature-specific terms.
- No API calls.
- Can emit events through props.

### Organisms

Reusable sections that combine multiple molecules.

Examples:

- `DataToolbar`
- `DataTable`
- `ModalForm`
- `AppNavbar`
- `AppSidebar`

Rules:

- Generic and configurable.
- No feature API imports.
- No hard-coded business module names.

### Templates

Reusable page-level shells.

Examples:

- `CrudPage`
- `AuthenticatedLayout`
- `SplitAuthLayout`

Rules:

- Layout and composition only.
- Feature data is passed through props or render functions.

## Shared Folder Contract

Use `shared` for code used across multiple features.

Allowed:

- Reusable UI.
- Generic hooks.
- API client and API error handling.
- Session storage helpers.
- Export helpers.
- Formatting helpers.
- Shared generic types.

Not allowed:

- Feature-specific API calls.
- Feature-specific table columns.
- Feature-specific form schemas.
- Feature-specific permission logic unless generalized.

## Import Rules

Allowed:

```text
app -> features
features -> shared
features -> lib
shared -> lib
```

Not allowed:

```text
shared -> features
features/users -> features/roles/internal-file
app -> features/users/components/users-table
app -> shared/api/client for route business logic
```

Cross-feature communication must use one of these:

- Shared API contracts.
- A public feature export.
- Backend-composed response data.

## Naming Rules

- Route folders: kebab-case.
- Feature folders: kebab-case.
- Component files: kebab-case.
- React components: PascalCase.
- Hooks: `useSomething`.
- Types: PascalCase.
- API functions: verb + noun, for example `listUsers`, `createUser`.
- Shared UI components should avoid business names.

## Current Feature Mapping

Current modules should be refactored into this target mapping:

```text
features/admin/components/admin-workspaces.tsx
```

Should be split into:

```text
features/tenants/
  index.tsx
  types.ts
  hooks.ts
  api.ts
  components/
    tenants-form.tsx
    tenants-table.tsx

features/users/
  index.tsx
  types.ts
  hooks.ts
  api.ts
  components/
    users-form.tsx
    users-table.tsx
    users-filters.tsx
    role-badges.tsx

features/roles/
  index.tsx
  types.ts
  hooks.ts
  api.ts
  components/
    roles-form.tsx
    roles-table.tsx
```

Current shared UI should be moved:

```text
features/shared/components/confirmation-dialog.tsx
```

To:

```text
shared/molecules/confirmation-dialog.tsx
```

Current API client/session should be moved:

```text
lib/api/client.ts -> shared/api/client.ts
features/auth/session.ts -> shared/auth/session.ts
```

`lib` remains for framework infrastructure such as environment and query provider.

## CRUD Feature Standard

Every CRUD feature should provide:

- Route page under `app/(app)/{feature}/page.tsx`.
- Feature view under `features/{feature}/index.tsx`.
- Typed API wrapper.
- Feature hooks for list, create, update, delete.
- Search.
- Filters.
- Excel export.
- PDF print.
- Pagination with page size.
- Loading state.
- Empty state.
- Error state.
- Add/edit modal form.
- Confirmation for create, update, delete.
- Tenant-scoped behavior when applicable.
- Unauthorized behavior through shared API client.

## Scaffolding Target

The frontend scaffold command should generate:

```text
apps/web/app/(app)/{feature}/
  page.tsx
  loading.tsx
  error.tsx

apps/web/features/{feature}/
  index.tsx
  types.ts
  hooks.ts
  api.ts
  schema.ts
  components/
    {feature}-form.tsx
    {feature}-table.tsx
    {feature}-filters.tsx
```

Generated route page:

```tsx
import { ExampleView } from "@/features/example";

export default function ExamplePage() {
  return <ExampleView />;
}
```

Generated feature index:

```tsx
"use client";

export function ExampleView() {
  return null;
}
```

## Migration Plan

Refactor should be done in small safe steps:

1. Create `shared` atomic design folders.
2. Move generic UI components from existing feature folders into `shared`.
3. Move API client and session helpers into `shared`.
4. Split `admin` workspace into `tenants`, `users`, and `roles` features.
5. Convert each `app/(app)/*/page.tsx` into thin route composition only.
6. Add feature hooks and types per module.
7. Update imports.
8. Run lint, typecheck, build, and backend smoke tests.
9. Commit after each stable step.

## Acceptance Criteria

- `apps/web/app` contains only route composition files, layouts, loading, and error boundaries.
- Each route has a matching feature view.
- Each feature exposes a public view through `index.tsx`.
- Feature state and data logic are in `hooks.ts`.
- Feature types are in `types.ts`.
- Shared reusable UI follows atomic design.
- No shared component imports from feature folders.
- New feature creation does not require editing existing feature internals.
- Existing UX behavior remains unchanged after refactor.
- Lint, typecheck, and build pass.
