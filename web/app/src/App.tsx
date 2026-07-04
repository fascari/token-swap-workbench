import { Activity, ArrowRightLeft, Blocks, Bot, Power, RefreshCcw, Send } from "lucide-react";
import { useState } from "react";

import {
  type Block,
  type BotAction,
  type BotResponse,
  fetchBlocks,
  fetchQuote,
  fetchStatus,
  runBots,
  submitTransaction,
} from "./api";

const TOKENS = ["NEX", "ETH", "DOGE"];

const defaultAccountID = 2;
const defaultInputToken = "NEX";
const defaultOutputToken = "ETH";
const defaultAmountIn = 10;
const defaultBlockCount = 10;
const defaultBotAmount = 10;

type RequestState = "idle" | "success" | "error";
type BotMode = "create" | "stop";

function wait(ms: number): Promise<void> {
  return new Promise((resolve) => {
    window.setTimeout(resolve, ms);
  });
}

function describeBlocks(blocks: Block[]): string {
  const nonEmptyBlocks = blocks.filter((block) => block.transactions.length > 0);
  if (nonEmptyBlocks.length === 0) {
    return "Recent blocks loaded, but the current window contains no transactions yet. Load again or increase Count.";
  }

  const newestMatch = nonEmptyBlocks[0];
  return `Loaded ${blocks.length} recent blocks. Block #${newestMatch.id} contains ${newestMatch.transactions.length} tx.`;
}

function describeBotSummary(summary: BotResponse): string {
  if (summary.action === "stop" && summary.all) {
    return `Stopped all bots: ${summary.active_bots} active bots remain, ${summary.attempted_operations} attempted operations.`;
  }

  return `${summary.action} ${summary.requested_amount}: ${summary.active_bots} active bots, ${summary.attempted_operations} attempted operations.`;
}

export function App() {
  const [accountID, setAccountID] = useState(defaultAccountID);
  const [inputToken, setInputToken] = useState(defaultInputToken);
  const [outputToken, setOutputToken] = useState(defaultOutputToken);
  const [amountIn, setAmountIn] = useState(defaultAmountIn);
  const [blockCount, setBlockCount] = useState(defaultBlockCount);
  const [botMode, setBotMode] = useState<BotMode>("create");
  const [stopAllBots, setStopAllBots] = useState(false);
  const [botAmount, setBotAmount] = useState(defaultBotAmount);
  const [status, setStatus] = useState("unknown");
  const [lastAction, setLastAction] = useState("Ready.");
  const [transactionStatus, setTransactionStatus] = useState("not submitted");
  const [requestState, setRequestState] = useState<RequestState>("idle");
  const [estimatedOutput, setEstimatedOutput] = useState<number | null>(null);
  const [botSummary, setBotSummary] = useState<BotResponse | null>(null);
  const [blocks, setBlocks] = useState<Block[] | null>(null);

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
      setLastAction(describeBlocks(result));
    } catch (error) {
      setRequestState("error");
      setLastAction(error instanceof Error ? error.message : "Blocks request failed.");
    }
  }

  async function quoteTransaction() {
    try {
      const result = await fetchQuote(inputToken, outputToken, amountIn);
      setEstimatedOutput(result.amount_out);
      setRequestState("success");
      setLastAction(`${amountIn} ${inputToken} quotes to ${result.amount_out} ${outputToken}.`);
    } catch (error) {
      setRequestState("error");
      setLastAction(error instanceof Error ? error.message : "Quote request failed.");
    }
  }

  async function sendTransaction() {
    try {
      const result = await submitTransaction({
        account_id: accountID,
        in_token: inputToken,
        out_token: outputToken,
        amount_in: amountIn,
      });
      setTransactionStatus(result.status);
      setRequestState("success");
      setLastAction("Transaction submitted. Waiting for block inclusion...");
      await wait(1200);
      await loadRecentBlocks(blockCount);
    } catch (error) {
      setRequestState("error");
      setLastAction(error instanceof Error ? error.message : "Transaction request failed.");
    }
  }

  async function submitBotAction() {
    try {
      const action: BotAction = botMode === "create" ? "create" : "stop";
      const result = await runBots({
        action,
        amount: stopAllBots ? undefined : botAmount,
        all: action === "stop" && stopAllBots,
      });
      setBotSummary(result);
      setRequestState("success");
      setLastAction(describeBotSummary(result));
    } catch (error) {
      setRequestState("error");
      setLastAction(error instanceof Error ? error.message : "Bot request failed.");
    }
  }

  return (
    <main className="app-shell">
      <section className="app-header">
        <div>
          <p className="eyebrow">Nexus Chain Adapter</p>
          <h1>Token Transaction Workbench</h1>
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
          <span className="status-label">Transaction status</span>
          <strong>{transactionStatus}</strong>
        </div>
      </section>

      <section className="workspace-grid">
        <section className="panel">
          <header className="panel-header">
            <ArrowRightLeft size={18} />
            <h2>Quote and Transaction</h2>
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
              <select value={inputToken} onChange={(event) => setInputToken(event.target.value)}>
                {TOKENS.map((token) => (
                  <option key={token} value={token}>
                    {token}
                  </option>
                ))}
              </select>
            </label>
            <label>
              <span>To</span>
              <select value={outputToken} onChange={(event) => setOutputToken(event.target.value)}>
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
            <strong>{estimatedOutput === null ? "Not quoted yet" : estimatedOutput.toFixed(8)}</strong>
          </div>

          <div className="actions-row">
            <button className="primary-button" type="button" onClick={quoteTransaction}>
              <Activity size={18} />
              <span>Quote</span>
            </button>
            <button className="primary-button" type="button" onClick={sendTransaction}>
              <Send size={18} />
              <span>Submit Transaction</span>
            </button>
          </div>
        </section>

        <section className="panel">
          <header className="panel-header">
            <Bot size={18} />
            <h2>Bot Orchestrator</h2>
          </header>

          <div className="bot-control">
            <label>
              <span>Action</span>
              <select
                value={botMode}
                onChange={(event) => {
                  const nextMode = event.target.value as BotMode;
                  setBotMode(nextMode);
                  if (nextMode === "create") {
                    setStopAllBots(false);
                  }
                }}
              >
                <option value="create">Create</option>
                <option value="stop">Stop</option>
              </select>
            </label>
            <label>
              <span>Amount</span>
              <input
                disabled={botMode === "stop" && stopAllBots}
                max="100"
                min="1"
                type="number"
                value={botAmount}
                onChange={(event) => setBotAmount(Number(event.target.value))}
              />
            </label>
            <button className="primary-button" type="button" onClick={submitBotAction}>
              <Power size={18} />
              <span>{botMode === "create" ? "Create" : stopAllBots ? "Stop all" : "Stop"}</span>
            </button>
          </div>

          <label className="bot-all-toggle" data-disabled={botMode !== "stop"}>
            <input
              checked={stopAllBots}
              disabled={botMode !== "stop"}
              type="checkbox"
              onChange={(event) => setStopAllBots(event.target.checked)}
            />
            <span>All active bots</span>
          </label>

          <div className="bot-summary">
            <span>Active bots</span>
            <strong>{botSummary === null ? "0" : botSummary.active_bots}</strong>
            <span>Operations</span>
            <strong>
              {botSummary === null
                ? "0 accepted / 0 failed"
                : `${botSummary.accepted_operations} accepted / ${botSummary.failed_operations} failed`}
            </strong>
          </div>
        </section>

        <section className="panel blocks-panel">
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
            A small window can miss the block that contains the submitted transaction. Increase Count or load
            again after one or two seconds if needed.
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

export default App;
