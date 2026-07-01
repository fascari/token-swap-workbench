# Token Swap Workbench

Token Swap Workbench is a Go and React application that exposes a small HTTP
API and operator UI on top of a Rust chain service.

## Overview

The repository provides:

- a Go API for chain status, quotes, swap submission, and recent blocks
- a React frontend for manual validation of the integration flow
- a local development workflow driven by `mise`

## Quick Start

Start the Rust chain service in the sibling repository:

```bash
cd ../<rust-chain-repo>
mise run chain
```

Then start this project:

```bash
cp .env.example .env
npm install
mise run run
```

Open `http://localhost:8080/`.

## Project Structure

- `cmd/api`: application entrypoint
- `internal/app/chain`: handlers, use cases, and domain types for chain
  operations
- `internal/bootstrap`: server wiring and startup
- `web/app`: React frontend source
- `docs`: project documentation and API contract

## Documentation

- [Documentation index](docs/README.md)
- [Testing guide](docs/guides/testing.md)
