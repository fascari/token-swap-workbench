# Frontend Usage Guide

This guide explains how to use the workbench screen and how to interpret the
visible results.

## Access

The normal entry point is `http://localhost:8080/`.

When the Go API is started with `mise run run`, it already serves the frontend
at that URL. A separate frontend server is not required for normal usage.

The optional Vite mode at `http://localhost:5173/` exists only for frontend
development.

## Example States

### Initial state

Before any action is executed:

- chain status can still be `unknown`
- no quote has been requested yet
- no blocks have been loaded yet

### Quoted state

After a valid quote request:

- `Last action` shows the quoted conversion
- `Estimated output` displays the returned amount
- transaction status is still `not submitted`

### Post-transaction blocks state

After a successful transaction submission and a later blocks load:

- transaction status is `submitted`
- `Last action` reports which recent block window was loaded
- a block with `1 tx` indicates that a transaction was included in chain
  history

## Screen Areas

### Chain Status Area

This area reflects `GET /v1/chain/status`.

At the current implementation level, it should be read as a lightweight
availability check:

- `ok` means the Go API can still reach the Rust chain service
- it does not expose account balances or a richer chain snapshot

### Quote And Transaction Area

This area is used to build the transaction submitted by the current UI:

- `Account`: the account id used for the transaction
- `From`: input token
- `To`: output token
- `Amount`: input amount

`Quote` requests an estimated output amount. `Submit Transaction` sends the
transaction to the backend through `POST /v1/transactions`.

### Recent Blocks Area

This area displays recent blocks returned by `GET /v1/blocks?n=<count>`.

Each card should be read as one chain history entry:

- `Block #53` is the block identifier
- the date and time show when that block was produced
- `0 tx` means that block contains no included transactions
- `1 tx` means that block contains one included transaction

The blocks list is not a list of swaps. It is a list of history entries
produced by the Rust chain. Some entries are empty. Some contain one or more
transactions.

### Bot Orchestrator Area

This area controls the Task 2 bot manager through `POST /v1/bots`.

Controls:

- `Action`: selects `Create` or `Stop`.
- `Amount`: how many bots the command applies to (1 to 100). It is disabled
  when `All active bots` is checked.
- `All active bots`: available only with `Stop`; stops every active bot and
  ignores `Amount`.
- The action button relabels itself based on the current selection: `Create`
  in create mode, `Stop` in stop mode, and `Stop all` when `All active bots`
  is checked.

Readouts:

- `Active bots`: how many bots are currently running after the last command.
- `Operations`: accepted versus failed transaction submissions reported by the
  manager (`X accepted / Y failed`).

Active bots submit randomized send or swap transaction envelopes to the Rust
chain through its documented `POST /transaction` API.

## Manual Flow

1. Open `http://localhost:8080/`.
2. Click `Refresh`.
3. Confirm that chain status becomes `ok`.
4. Choose an account, token pair, and amount.
5. Click `Quote`.
6. Confirm that an estimated output appears.
7. Click `Submit Transaction`.
8. Confirm that transaction status becomes `submitted`.
9. Click `Load` in the recent blocks area.
10. If the newest blocks still show `0 tx`, wait one or two seconds and load
    blocks again.

## Bot Orchestrator Flow

1. Open `http://localhost:8080/`.
2. Set `Action` to `Create`.
3. Set `Amount` to `10`.
4. Click `Create`.
5. Confirm that `Active bots` becomes `10`.
6. Set `Action` to `Stop`.
7. Set `Amount` to `5`.
8. Click `Stop`.
9. Confirm that `Active bots` becomes `5`.
10. Check `All active bots` to stop the remaining bots.
11. Click `Stop all`.
12. Confirm that `Active bots` becomes `0`.
13. Load recent blocks to inspect chain activity produced by the bots.

## Post-Transaction Behavior In The Example UI

The current example UI applies a small amount of automation after `Submit
Transaction`:

- the transaction status changes to `submitted`
- the screen waits briefly for block production to advance
- the recent blocks list is loaded automatically

The blocks area uses `Count = 10` by default because a very small window can
miss the block that contains the submitted transaction.

If the blocks area still shows only `0 tx` entries after submission, that does
not mean the transaction failed. It usually means one of these cases:

- the transaction was accepted, but the currently loaded window does not yet include
  the block where it was recorded
- the block was produced slightly later than the current refresh moment
- the block containing the transaction already moved outside a very small recent-block
  window

The practical next step is:

1. wait one or two seconds
2. click `Load` again
3. if needed, increase `Count`

The most useful confirmation currently available in the example UI is a later
block showing a non-zero transaction count such as `1 tx`.

## What To Verify In The UI

The current UI is useful for confirming:

- the Go API is reachable
- the Rust chain service is reachable through the Go API
- quote calculation is working
- a transaction request was accepted
- bot create/stop commands are accepted by the API
- the chain is still producing blocks
- a later block eventually shows a non-zero transaction count

After `Submit Transaction`, the most useful visual confirmation currently available in
the UI is a later block with `1 tx` or another non-zero count.

The current UI does not directly verify:

- account balance changes after the transaction
- the raw transaction payload shown by the Rust service

Deeper confirmation of the underlying block payload belongs to the API response
from `GET /v1/blocks` and is described in
[System Operation Flow](../system/operation-flow.md).
