# Development Guide

## Prerequisites

- Go 1.26.1+
- Docker & Docker Compose
## First Time Setup

1. **Clone the repo**
```bash
git clone <repo-url>
cd token-swap-workbench
```

2. **Copy environment file**
```bash
cp .env.example .env
```

3. **Install dependencies**
```bash
go mod download
```


## Running Locally

Start the service:
```bash
mise run run
```

The API will be available at `http://localhost:8080`.

## Docker Compose Layout

No docker-compose.yml included. Set `include_dockerfile = true` when scaffolding.
## Common Workflows

### Adding a New Feature

1. Create domain directory: `internal/app/users/`
2. Add models: `models.go`
3. Add repository interface and implementation: `repository.go`
4. Add service: `service.go`
5. Add HTTP handler: `handler.go`
6. Wire up in `internal/bootstrap/bootstrap.go`
7. Add routes in `cmd/api/main.go`

### Database Changes

Database not included. Set `include_db = true` when scaffolding.
### Running Tests

```bash
mise run test        # Unit tests
mise run e2e         # E2E tests
mise run lint        # Linter
```

## Mise Quick Reference

| Task | Description |
|------|-------------|
| `mise run test` | Run unit tests |
| `mise run lint` | Run golangci-lint |
| `mise run tidy` | Tidy dependencies |
| `mise run build` | Build binary |
| `mise run fmt` | Format code |
| `mise run run` | Run service locally |
| `mise run e2e` | Run e2e tests |
| `mise run e2e:up` | Start e2e postgres |
| `mise run e2e:down` | Stop e2e postgres |
| `mise run mocks` | Generate mocks |
