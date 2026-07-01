# Development Guide

## Prerequisites

- Go `1.26.1+`
- Node `22+`
- the Rust chain service repository available as a sibling directory

## Setup

```bash
cp .env.example .env
go mod download
npm install
```

## Configuration

The main local settings are:

- `HTTP_PORT=8080`
- `CHAIN_BASE_URL=http://127.0.0.1:3000`
- `LOG_LEVEL=info`
- `LOG_FORMAT=json`

## Local Run

Start the Rust chain service first:

```bash
cd ../<rust-chain-repo>
mise run chain
```

Then, in this repository, start the integrated runtime:

```bash
mise run run
```

This command builds the React frontend into `web/static` and starts the Go API.
Once it is running, the application is available directly at
`http://localhost:8080/`. No separate frontend server is required for normal
usage.

## Frontend Development Mode

If changes are being made to the React UI, keep the Go API running on `:8080`
and start Vite in another terminal:

```bash
mise run frontend:dev
```

The frontend is then available at `http://localhost:5173/`. In this mode, Vite
proxies `/v1` and `/health` requests back to the Go API at
`http://127.0.0.1:8080`.

## Smoke Check

The fastest manual verification path is:

1. open `http://localhost:8080/`
2. click `Refresh`
3. request a quote
4. submit a swap
5. load recent blocks

Detailed UI behavior and manual flow are documented in
[Frontend Usage Guide](../frontend/usage-guide.md). System-level request flow is
documented in [System Overview](../system/README.md). The API contract lives in
[openapi.yaml](../openapi.yaml).

## Common Commands

| Command | Purpose |
| --- | --- |
| `mise run run` | Build frontend assets and run the Go API |
| `mise run frontend:dev` | Start the Vite development server |
| `mise run build` | Build the Go API binary into `build/` |
| `mise run test` | Run unit tests |
| `mise run integration` | Run deterministic integration tests |
| `mise run lint` | Run `golangci-lint` |
| `mise run fmt` | Run `go fmt ./...` |
| `mise run mocks` | Regenerate mocks |
