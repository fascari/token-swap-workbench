import { RefreshCcw } from "lucide-react";

type RequestState = "idle" | "success" | "error";

type StatusBarProps = {
  status: string;
  lastAction: string;
  transactionStatus: string;
  requestState: RequestState;
  onRefresh: () => void;
};

export function StatusBar({ status, lastAction, transactionStatus, requestState, onRefresh }: StatusBarProps) {
  return (
    <>
      <section className="app-header">
        <div>
          <p className="eyebrow">Nexus Chain Adapter</p>
          <h1>Token Transaction Workbench</h1>
        </div>
        <button className="icon-button" type="button" onClick={onRefresh}>
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
    </>
  );
}
