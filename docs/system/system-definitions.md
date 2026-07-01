# System Definitions

This document defines the main concepts used by the application in neutral
product terms. It is intended for readers without blockchain background.

## Core Concepts

### Blockchain

A blockchain is an append-only sequence of records. Each record references the
previous one, which creates an ordered history that is difficult to rewrite
without rebuilding the chain from that point forward.

### Block

A block is one entry in that sequence. It works like a small batch of chain
history created at a specific moment in time.

In practical terms, a block in this project contains:

- a block identifier
- a timestamp
- zero or more transactions included in that block

The Rust chain service produces blocks continuously, even when no new
transaction is submitted. That is why some blocks appear with
`transactions: []`: the chain created a new history entry, but there was no
transaction to record in that specific block.

When a swap is submitted, it does not become visible as a "block" by itself.
Instead, the swap transaction is included inside a later block. In the UI, the
recent blocks panel should be read as a rolling history feed: each row is one
moment of chain history, and the transaction list inside that row shows what
happened at that moment.

### Transaction

A transaction is an operation submitted to the chain. In this workbench, the
swap submission creates a transaction that is later visible through the latest
blocks.

### Account

An account is an identity that can hold token balances. The chain status
runtime maintains a small set of sample accounts that can be used in swaps and
other transactions.

### Token

A token is a unit of value tracked by the chain. The service exposes sample
tokens such as `NEX`, `ETH`, and `DOGE`. A balance such as `10 NEX` means that
the account currently holds ten units of the `NEX` token.

### Swap

A swap exchanges one token for another. A request such as "swap `10 NEX` into
`ETH`" means that the caller wants to spend `10` units of `NEX` and receive the
equivalent amount of `ETH` according to the quote rules exposed by the chain
service.

### Quote

A quote is a calculation returned before the swap is submitted. It estimates
the output amount for a given input amount and token pair. It is similar to a
price preview in a traditional exchange flow. Because the Rust service uses an
AMM-style pool, the returned value can change after swaps alter pool reserves.

## API-Facing Terms

### Chain Status

Chain status is the lightweight response returned by `GET /v1/chain/status`.
In the current implementation it acts as an upstream availability check and
returns a simple service status instead of a full account-and-balance snapshot.

### `amount_in`

`amount_in` is the input quantity sent to the quote or swap endpoints. It
represents how much of the source token should be spent.

### Recent Blocks

Recent blocks are returned by `GET /v1/blocks?n=<count>`. They provide a simple
way to inspect the latest chain activity after a swap is submitted. Because the
Rust chain produces blocks continuously, the submitted swap may appear in a
later block rather than in the first block fetched immediately after the
request returns.
