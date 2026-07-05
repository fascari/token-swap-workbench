import { useState } from "react";

import type { BotAction, BotMode, BotResponse } from "./api/types";
import { fetchQuote, fetchStatus, runBots, submitTransaction } from "./api/client";
import { BlocksPanel } from "./components/BlocksPanel";
import { BotPanel } from "./components/BotPanel";
import { QuotePanel } from "./components/QuotePanel";
import { StatusBar } from "./components/StatusBar";
import { useBlocks } from "./hooks/useBlocks";
import { usePolling } from "./hooks/usePolling";

const TOKENS = ["NEX", "ETH", "DOGE"];

const defaultAccountID = 2;
const defaultInputToken = "NEX";
const defaultOutputToken = "ETH";
const defaultAmountIn = 10;
const defaultBotAmount = 10;

type RequestState = "idle" | "success" | "error";

function wait(ms: number): Promise<void> {
  return new Promise((resolve) => {
    window.setTimeout(resolve, ms);
  });
}

function describeBlocks(blocks: { id: number; transactions: unknown[] }[]): string {
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
  const [botMode, setBotMode] = useState<BotMode>("create");
  const [stopAllBots, setStopAllBots] = useState(false);
  const [botAmount, setBotAmount] = useState(defaultBotAmount);
  const [status, setStatus] = useState("unknown");
  const [lastAction, setLastAction] = useState("Ready.");
  const [transactionStatus, setTransactionStatus] = useState("not submitted");
  const [requestState, setRequestState] = useState<RequestState>("idle");
  const [estimatedOutput, setEstimatedOutput] = useState<number | null>(null);
  const [botSummary, setBotSummary] = useState<BotResponse | null>(null);

  const { blocks, blockCount, setBlockCount, load: loadRecentBlocks } = useBlocks();
  const { isPolling, start: startPolling, stop: stopPolling } = usePolling();

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
      await loadRecentBlocks();
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

  function handleStartPolling() {
    startPolling(loadRecentBlocks, 2000);
  }

  return (
    <main className="app-shell">
      <StatusBar
        status={status}
        lastAction={lastAction}
        transactionStatus={transactionStatus}
        requestState={requestState}
        onRefresh={refreshStatus}
      />

      <section className="workspace-grid">
        <QuotePanel
          tokens={TOKENS}
          accountID={accountID}
          inputToken={inputToken}
          outputToken={outputToken}
          amountIn={amountIn}
          estimatedOutput={estimatedOutput}
          onAccountIDChange={setAccountID}
          onInputTokenChange={setInputToken}
          onOutputTokenChange={setOutputToken}
          onAmountInChange={setAmountIn}
          onQuote={quoteTransaction}
          onSubmit={sendTransaction}
        />

        <BotPanel
          botMode={botMode}
          botAmount={botAmount}
          stopAllBots={stopAllBots}
          botSummary={botSummary}
          onBotModeChange={setBotMode}
          onBotAmountChange={setBotAmount}
          onStopAllBotsChange={setStopAllBots}
          onSubmitBotAction={submitBotAction}
        />

        <BlocksPanel
          blocks={blocks}
          blockCount={blockCount}
          isPolling={isPolling}
          onBlockCountChange={setBlockCount}
          onLoad={() => loadRecentBlocks()}
          onStartPolling={handleStartPolling}
          onStopPolling={stopPolling}
        />
      </section>
    </main>
  );
}

export default App;
