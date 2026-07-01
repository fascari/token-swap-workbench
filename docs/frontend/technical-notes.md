# Frontend Technical Notes

## Runtime Topology

```text
React UI -> Go API -> Rust chain service
```

## Source Layout

```text
web/app/
  src/
    App.tsx
    api.ts
    main.tsx
    styles.css
```

`App.tsx` owns the screen state and action flow. `api.ts` contains the HTTP
calls to the Go API.

## Development and Build Modes

`mise run frontend:dev` starts Vite on `http://localhost:5173/` and proxies
`/health` and `/v1` to the Go API at `http://127.0.0.1:8080`.

`mise run run` builds `web/app` into `web/static`, then starts the Go API,
which serves:

- `/` from `web/static/index.html`
- `/static/*` for generated assets

The `/static/*` path is an implementation detail of the production asset build.
The human-facing entry point remains `/`.

## Why React Is Kept Small

The frontend is intentionally simple:

- one screen
- no client-side routing
- local component state only
- direct API calls without a frontend BFF layer

That keeps the integration path easy to validate while still leaving room for a
future expansion if a richer operator workflow is needed.
