# Operation Flow

This document explains what each visible operation triggers in the system and
how it maps to the Rust chain service.

## Runtime Topology

```text
React UI -> Go API -> Rust chain service
```

The React application collects user input, the Go API exposes a stable HTTP
surface, and the Rust service performs the chain-specific work.

## Refresh Status

The `Refresh` action reloads the current chain view.

1. The frontend calls `GET /v1/chain/status`.
2. The Go API forwards the request to the Rust service through the chain
   client.
3. The response returns a lightweight service status.
4. The frontend updates the status indicator.

From a product perspective, this operation answers: "What assets exist right
now, and what can be swapped?" In the current implementation the answer is only
partial because this endpoint behaves as a connectivity check.

## Request Quote

The `Get Quote` action previews the expected output for a token conversion.

1. The operator selects the source token, target token, and input amount.
2. The frontend calls `GET /v1/quote`.
3. The Go API validates the payload and forwards the quote request to the Rust
   service.
4. The Rust service calculates the expected output and returns the quote.
5. The frontend renders the returned amount and keeps the payload available for
   the next step.

From a product perspective, this operation answers: "If this trade is executed
now, what should come back?" Because the upstream pool state changes after
swaps, repeating the same quote later can return a different `amount_out`.

## Submit Swap

The `Submit Swap` action executes the token exchange.

1. The frontend sends `POST /v1/transactions` with the selected account, token pair,
   and input amount.
2. The Go API forwards the request to the Rust service.
3. The Rust service accepts the swap transaction and returns a submission
   status.
4. The frontend shows the result and the operator can inspect recent blocks to
   confirm when the transaction is included by the block producer.

From a product perspective, this operation answers: "Execute the exchange and
return the chain result."

### What To Verify After Submit Swap

When the UI returns `submitted`, the operator has verified only that:

- the Go API accepted the request
- the Go API could forward the request to the Rust chain service
- the Rust chain service accepted the swap submission

At this stage, the operator has not yet verified block inclusion.

## Load Recent Blocks

The `Load Blocks` action fetches the latest chain history.

1. The frontend sends `GET /v1/blocks?n=<count>`.
2. The Go API requests the latest blocks from the Rust service.
3. The response returns the most recent block list.
4. The frontend renders the history panel.
5. If the most recent blocks still show empty transaction arrays, another fetch
   a second or two later can reveal the block where the submitted swap was
   included.

From a product perspective, this operation answers: "What happened on the chain
most recently, and is the submitted swap visible in that history?"

### What To Verify In The Blocks View

After a swap submission, the blocks view is useful for verifying three things:

1. The chain is still producing new blocks.
2. A later block may include a transaction entry instead of `transactions: []`.
3. In the current UI, that inclusion appears as a block with a non-zero
   transaction count such as `1 tx`.

This means the operator should not expect the newest block fetched immediately
after submission to always contain the swap. A valid manual check is:

1. submit the swap and confirm the UI returns `submitted`
2. load recent blocks
3. if the newest blocks are still empty, wait one or two seconds
4. load recent blocks again
5. confirm that a later block eventually shows a non-zero transaction count
6. if deeper confirmation is needed, inspect the raw API response for
   `GET /v1/blocks`, where the Rust transaction envelope appears as
   `{"Swap": {...}}`

What the current UI does not verify directly is balance change by account. The
screen is currently better suited to confirm request acceptance, quote
calculation, ongoing block production, and eventual swap inclusion in block
history.
