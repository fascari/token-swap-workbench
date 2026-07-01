# Evolution Notes

This document records architectural extensions that would fit the current
service if the scope grows beyond direct request forwarding.

## Off-Chain Persistence

The current implementation does not persist application state outside the Rust
chain service. A persistence layer on the Go side starts to make sense when the
system needs any of the following:

- operator audit history
- search and filtering over past swaps
- request deduplication or idempotency tracking
- reporting or analytics disconnected from the chain runtime
- workflow state that should survive a Rust service restart

At that point, a relational store such as PostgreSQL is a reasonable default.

## Event-Driven Extensions

An event-driven model becomes useful when the swap flow needs asynchronous side
effects such as notifications, reconciliation jobs, or downstream projections.
Examples include:

- publishing swap-completed events to another service
- building read models optimized for dashboards
- decoupling submission from long-running enrichment steps

The current synchronous API is sufficient while each request maps directly to a
single upstream chain call.

## Dual Write and Outbox

If a future version both stores data locally and publishes events, the design
must avoid dual-write inconsistency. A common approach is the outbox pattern:

1. write the business record and outbox event in the same database transaction
2. publish the outbox event asynchronously
3. mark the event as processed after a successful publish

That pattern keeps the local source of truth and emitted events aligned.

## Idempotency and Retries

If clients or background workers can retry swap submission, idempotency becomes
important. The Go layer would then need stable request identifiers and a way to
detect whether the same logical operation has already been accepted.

## Why This Still Matters in the Current Design

Even though the current workbench is intentionally small, these topics are
natural follow-up questions because the service already sits between a user
interface and a chain runtime. The current structure leaves room to introduce
storage, asynchronous workflows, and richer orchestration without changing the
public API shape abruptly.
