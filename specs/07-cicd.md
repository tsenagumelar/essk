# ESSK Specs - CI/CD

## Platform

Use GitHub Actions.

Required workflows:

```text
.github/workflows/ci.yml
.github/workflows/docker.yml
.github/workflows/release.yml
```

## CI Workflow

Trigger:

- Pull request to main.
- Push to main.

Jobs:

- Pre-commit policy checks.
- Backend lint.
- Backend test.
- Backend build.
- Frontend lint.
- Frontend typecheck.
- Frontend test.
- Dependency vulnerability scan.
- License policy scan.
- Docker Compose config validation.

Backend commands:

```text
pre-commit run --all-files
essk policy check
go test ./...
go vet ./...
golangci-lint run
govulncheck ./...
go build ./services/api-gateway
```

Frontend commands:

```text
pnpm install --frozen-lockfile
pnpm lint
pnpm typecheck
pnpm test
pnpm build
pnpm audit
```

Governance commands:

```text
essk policy check --all
essk dependency check
essk license check
```

## Docker Workflow

Trigger:

- Push to main.
- Version tag.

Build images:

- `essk-backend`.
- `essk-web`.
- `essk-admin` when admin app exists.

Use:

- Docker Buildx.
- GitHub Container Registry by default.
- Cache from registry.

Tags:

- `latest` for main only.
- Git SHA.
- Semantic version tag.

## Release Workflow

Trigger:

- Git tag `v*`.

Responsibilities:

- Run full CI.
- Build and push versioned Docker images.
- Generate changelog.
- Create GitHub release.
- Attach release SBOM when available.
- Publish rollback notes.

## Required Branch Policy

For `main`:

- Require pull request.
- Require CI passing.
- Require at least one review.
- Disallow force push.

## Secrets

GitHub Actions secrets:

- `GHCR_TOKEN` if not using default GitHub token.
- Deployment credentials later.

Never store application secrets in workflow files.

## Deployment Strategy

Phase 1:

- CI validates code.
- Docker images are buildable.
- Deployment remains manual/local.

Phase 2:

- Add staging deployment workflow.
- Run migrations as a controlled job.
- Add smoke tests after deployment.

Production future:

- Blue-green or rolling deployment.
- Database migration review step.
- Rollback documentation.
- Smoke tests after deployment.
- Automatic rollback trigger for failed smoke tests where supported.
