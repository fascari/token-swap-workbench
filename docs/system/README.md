# System Overview

This section groups the documents that explain the domain model, the request
flow, and the main evolution paths for the service.

## Documents

- [System Definitions](system-definitions.md): core terms used by the service
  and the API.
- [Operation Flow](operation-flow.md): what each UI action triggers and how the
  Go service interacts with the Rust chain service.
- [Evolution Notes](evolution-notes.md): architectural extensions that fit the
  current design if the system grows beyond the current scope.

## Scope

The current implementation keeps the Rust chain service as the execution
engine. The Go application acts as an HTTP adapter and workflow-oriented API
layer. The React frontend exists to validate the integration path and expose
the chain operations in a compact interface.
