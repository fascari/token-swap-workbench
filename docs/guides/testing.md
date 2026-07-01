# Testing Guide

The project uses two kinds of automated tests:

- Unit tests. Fast, isolated tests that mock the collaborators of a single
  package (handlers, use cases, the chain client).
- Integration tests. End-to-end tests that drive the real HTTP stack and
  replace only the external service with an in-memory stub.

Run everything with:

```bash
mise run test
```

That task runs the unit and integration suites together. The two can also be
run on their own with `mise run unit` and `mise run integration`.

## Unit Tests

Run them with:

```bash
mise run unit
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

They drive the whole HTTP stack for one operation at a time, with the external
service swapped for a controllable stub, so every run is deterministic and needs
no database or external network.

See the [Integration testing guide](integration-tests.md) for the full picture,
including what runs for real, who starts the server, how a request flows through
the stack, why the input and output fixtures differ, and how the suite and
fixtures are organized.
