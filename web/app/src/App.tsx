import { useState } from "react";
import {
  Activity,
  ArrowRightLeft,
  Blocks,
  RefreshCcw,
  Send,
} from "lucide-react";

import {
  fetchBlocks,
  fetchQuote,
  fetchStatus,
  submitSwap,
  type BlockResponse,
} from "./api";

const TOKENS = ["NEX", "ETH", "DOGE"] as const;
const defaultAccountID = 2;
const defaultInputToken = "NEX";
const defaultOutputToken = "ETH";
const defaultAmountIn = 10;
const defaultBlockCount = 10;

type RequestState = "idle" | "success" | "error";

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => {
    window.setTimeout(resolve, ms);
  });
}

function blocksSummary(blocks: BlockResponse[]): string {
  const nonEmptyBlocks = blocks.filter((block) => block.transactions.length > 0);

  if (nonEmptyBlocks.length === 0) {
    return "Recent blocks loaded, but the current window contains no transactions yet. Load again or increase Count.";
  }

  const newestMatch = nonEmptyBlocks[0];

  return `Loaded ${blocks.length} recent blocks. Block #${newestMatch.id} contains ${newestMatch.transactions.length} tx.`;
}

export default function App() {
  const [accountID, setAccountID] = useState(defaultAccountID);
  const [inputToken, setInputToken] = useState(defaultInputToken);
  const [outputToken, setOutputToken] = useState(defaultOutputToken);
  const [amountIn, setAmountIn] = useState(defaultAmountIn);
  const [blockCount, setBlockCount] = useState(defaultBlockCount);
  const [status, setStatus] = useState("unknown");
  const [lastAction, setLastAction] = useState("Ready.");
  const [swapStatus, setSwapStatus] = useState("not submitted");
  const [requestState, setRequestState] = useState<RequestState>("idle");
  const [estimatedOutput, setEstimatedOutput] = useState<number | null>(null);
  const [blocks, setBlocks] = useState<BlockResponse[] | null>(null);

  async function refreshStatus() {
    try {
      const result = await fetchStatus();
      setStatus(result.status);
      setRequestState("success");
      setLastAction("Chain status refreshed.");
    } catch (error) {
      setRequestState("error");
      setLastAction(error instanceof Error ? error.message : "Status request failed.");
    }
  }

  async function loadRecentBlocks(count = blockCount) {
    try {
      const result = await fetchBlocks(count);
      setBlocks(result);
      setRequestState("success");
      setLastAction(blocksSummary(result));
    } catch (error) {
      setRequestState("error");
      setLastAction(error instanceof Error ? error.message : "Blocks request failed.");
    }
  }

  async function requestQuote() {
    try {
      const result = await fetchQuote(inputToken, outputToken, amountIn);
      setEstimatedOutput(result.amount_out);
      setRequestState("success");
      setLastAction(
        `${amountIn} ${inputToken} quotes to ${result.amount_out} ${outputToken}.`,
      );
    } catch (error) {
      setRequestState("error");
      setLastAction(error instanceof Error ? error.message : "Quote request failed.");
    }
  }

  async function sendSwap() {
    try {
      const result = await submitSwap({
        account_id: accountID,
        in_token: inputToken,
        out_token: outputToken,
        amount_in: amountIn,
      });

      setSwapStatus(result.status);
      setRequestState("success");
      setLastAction("Swap transaction submitted. Waiting for block inclusion...");

      await sleep(1200);
      await loadRecentBlocks(blockCount);
    } catch (error) {
      setRequestState("error");
      setLastAction(error instanceof Error ? error.message : "Swap request failed.");
    }
  }

  return (
    <main className="app-shell">
      <section className="app-header">
        <div>
          <p className="eyebrow">Nexus Chain Adapter</p>
          <h1>Token Swap Workbench</h1>
        </div>

        <button className="icon-button" type="button" onClick={refreshStatus}>
          <RefreshCcw size={18} />
          <span>Refresh</span>
        </button>
      </section>

      <section className="status-strip" data-state={requestState}>
        <div className="status-card">
          <span className="status-label">Chain status</span>
          <strong>{status}</strong>
        </div>

        <div className="status-card status-card-wide">
          <span className="status-label">Last action</span>
          <strong>{lastAction}</strong>
        </div>

        <div className="status-card">
          <span className="status-label">Swap status</span>
          <strong>{swapStatus}</strong>
        </div>
      </section>

      <section className="workspace-grid">
        <section className="panel">
          <header className="panel-header">
            <ArrowRightLeft size={18} />
            <h2>Quote and Swap</h2>
          </header>

          <div className="field-grid">
            <label>
              <span>Account</span>
              <input
                min="1"
                type="number"
                value={accountID}
                onChange={(event) => setAccountID(Number(event.target.value))}
              />
            </label>

            <label>
              <span>From</span>
              <select
                value={inputToken}
                onChange={(event) => setInputToken(event.target.value)}
              >
                {TOKENS.map((token) => (
                  <option key={token} value={token}>
                    {token}
                  </option>
                ))}
              </select>
            </label>

            <label>
              <span>To</span>
              <select
                value={outputToken}
                onChange={(event) => setOutputToken(event.target.value)}
              >
                {TOKENS.map((token) => (
                  <option key={token} value={token}>
                    {token}
                  </option>
                ))}
              </select>
            </label>

            <label>
              <span>Amount</span>
              <input
                min="0.000001"
                step="0.000001"
                type="number"
                value={amountIn}
                onChange={(event) => setAmountIn(Number(event.target.value))}
              />
            </label>
          </div>

          <div className="quote-output">
            <span>Estimated output</span>
            <strong>
              {estimatedOutput === null ? "Not quoted yet" : estimatedOutput.toFixed(8)}
            </strong>
          </div>

          <div className="actions-row">
            <button className="primary-button" type="button" onClick={requestQuote}>
              <Activity size={18} />
              <span>Quote</span>
            </button>

            <button className="primary-button" type="button" onClick={sendSwap}>
              <Send size={18} />
              <span>Submit Swap</span>
            </button>
          </div>
        </section>

        <section className="panel">
          <header className="panel-header">
            <Blocks size={18} />
            <h2>Recent Blocks</h2>
          </header>

          <div className="blocks-control">
            <label>
              <span>Count</span>
              <input
                max="20"
                min="1"
                type="number"
                value={blockCount}
                onChange={(event) => setBlockCount(Number(event.target.value))}
              />
            </label>

            <button className="primary-button" type="button" onClick={() => loadRecentBlocks()}>
              <RefreshCcw size={18} />
              <span>Load</span>
            </button>
          </div>

          <p className="blocks-hint">
            A small window can miss the block that contains the submitted swap.
            Increase Count or load again after one or two seconds if needed.
          </p>

          <div className="blocks-list">
            {blocks === null ? (
              <p className="empty-state">No blocks loaded.</p>
            ) : (
              blocks.map((block) => (
                <article className="block-row" key={block.id}>
                  <div>
                    <span>{`Block #${block.id}`}</span>
                    <strong>{`${block.transactions.length} tx`}</strong>
                  </div>
                  <time>{new Date(block.timestamp * 1000).toLocaleString()}</time>
                </article>
              ))
            )}
          </div>
        </section>
      </section>
    </main>
  );
}
