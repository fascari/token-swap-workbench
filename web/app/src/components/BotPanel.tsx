import { Bot, Power } from "lucide-react";

import type { BotMode } from "../api/types";

type BotPanelProps = {
  botMode: BotMode;
  botAmount: number;
  stopAllBots: boolean;
  botSummary: { active_bots: number; accepted_operations: number; failed_operations: number } | null;
  onBotModeChange: (mode: BotMode) => void;
  onBotAmountChange: (amount: number) => void;
  onStopAllBotsChange: (all: boolean) => void;
  onSubmitBotAction: () => void;
};

export function BotPanel({
  botMode,
  botAmount,
  stopAllBots,
  botSummary,
  onBotModeChange,
  onBotAmountChange,
  onStopAllBotsChange,
  onSubmitBotAction,
}: BotPanelProps) {
  const actionLabel = botMode === "create" ? "Create" : stopAllBots ? "Stop all" : "Stop";

  return (
    <section className="panel">
      <header className="panel-header">
        <Bot size={18} />
        <h2>Bot Orchestrator</h2>
      </header>

      <div className="toolbar">
        <label className="form-group form-group--sm">
          <span>Action</span>
          <select
            className="select"
            value={botMode}
            onChange={(event) => {
              const nextMode = event.target.value as BotMode;
              onBotModeChange(nextMode);
              if (nextMode === "create") {
                onStopAllBotsChange(false);
              }
            }}
          >
            <option value="create">Create</option>
            <option value="stop">Stop</option>
          </select>
        </label>

        <label className="form-group form-group--xs">
          <span>Amount</span>
          <input
            className="input"
            disabled={botMode === "stop" && stopAllBots}
            max="100"
            min="1"
            type="number"
            value={botAmount}
            onChange={(event) => onBotAmountChange(Number(event.target.value))}
          />
        </label>

        <label className="checkbox" data-disabled={botMode !== "stop"}>
          <input
            checked={stopAllBots}
            disabled={botMode !== "stop"}
            type="checkbox"
            onChange={(event) => onStopAllBotsChange(event.target.checked)}
          />
          <span>All active bots</span>
        </label>

        <button className="primary-button" type="button" onClick={onSubmitBotAction}>
          <Power size={16} />
          <span>{actionLabel}</span>
        </button>
      </div>

      <div className="stats-grid">
        <div className="stat-card">
          <span>Active bots</span>
          <strong>{botSummary === null ? "0" : botSummary.active_bots}</strong>
        </div>
        <div className="stat-card">
          <span>Operations</span>
          <strong>
            {botSummary === null
              ? "0 accepted / 0 failed"
              : `${botSummary.accepted_operations} accepted / ${botSummary.failed_operations} failed`}
          </strong>
        </div>
      </div>
    </section>
  );
}
