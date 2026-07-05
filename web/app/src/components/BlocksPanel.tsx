import { Blocks, Play, RefreshCcw, Square } from "lucide-react";

import type { Block } from "../api/types";

type BlocksPanelProps = {
  blocks: Block[] | null;
  blockCount: number;
  isPolling: boolean;
  onBlockCountChange: (count: number) => void;
  onLoad: () => void;
  onStartPolling: () => void;
  onStopPolling: () => void;
};

function BlockCard({ block }: { block: Block }) {
  return (
    <article className="block-row">
      <div>
        <span>{`Block #${block.id}`}</span>
        <strong>{`${block.transactions.length} tx`}</strong>
      </div>
      <time>{new Date(block.timestamp * 1000).toLocaleString()}</time>
    </article>
  );
}

export function BlocksPanel({
  blocks,
  blockCount,
  isPolling,
  onBlockCountChange,
  onLoad,
  onStartPolling,
  onStopPolling,
}: BlocksPanelProps) {
  return (
    <section className="panel blocks-panel">
      <header className="panel-header">
        <Blocks size={18} />
        <h2>Recent Blocks</h2>
      </header>

      <div className="toolbar">
        <label className="form-group form-group--xs">
          <span>Count</span>
          <input
            className="input"
            max="20"
            min="1"
            type="number"
            value={blockCount}
            onChange={(event) => onBlockCountChange(Number(event.target.value))}
          />
        </label>

        <button className="primary-button" type="button" onClick={onLoad}>
          <RefreshCcw size={16} />
          <span>Load</span>
        </button>

        <div className="toolbar-spacer" />

        <div className="poll-control">
          {isPolling ? (
            <button className="primary-button poll-active" type="button" onClick={onStopPolling}>
              <Square size={16} />
              <span>Stop Polling</span>
            </button>
          ) : (
            <button className="primary-button" type="button" onClick={onStartPolling}>
              <Play size={16} />
              <span>Start Polling</span>
            </button>
          )}
          {isPolling && <span className="poll-indicator" aria-label="Polling active" />}
        </div>
      </div>

      <div className="blocks-list">
        {blocks === null ? (
          <p className="empty-state">No blocks loaded.</p>
        ) : (
          blocks.map((block) => <BlockCard key={block.id} block={block} />)
        )}
      </div>
    </section>
  );
}
