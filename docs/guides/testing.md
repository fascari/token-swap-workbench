# Testing Guide

The project uses unit tests and deterministic integration tests.

## Unit Tests

Run the default suite with:

```bash
mise run test
```

Regenerate mocks with:

```bash
mise run mocks
```

## Integration Tests

Integration tests use the `integration` build tag and run with:

```bash
mise run integration
```

The generic integration suite starts an HTTP handler provided by each test
package and exposes `Expect()` for HTTP assertions. It does not know about
domain endpoints, fixtures, upstream services, or application modules.

Each integration package owns its own wiring, fixtures, and test doubles.
