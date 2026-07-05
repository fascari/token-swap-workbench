import { Activity, ArrowRightLeft, Send } from "lucide-react";

type QuotePanelProps = {
  tokens: string[];
  accountID: number;
  inputToken: string;
  outputToken: string;
  amountIn: number;
  estimatedOutput: number | null;
  onAccountIDChange: (id: number) => void;
  onInputTokenChange: (token: string) => void;
  onOutputTokenChange: (token: string) => void;
  onAmountInChange: (amount: number) => void;
  onQuote: () => void;
  onSubmit: () => void;
};

export function QuotePanel({
  tokens,
  accountID,
  inputToken,
  outputToken,
  amountIn,
  estimatedOutput,
  onAccountIDChange,
  onInputTokenChange,
  onOutputTokenChange,
  onAmountInChange,
  onQuote,
  onSubmit,
}: QuotePanelProps) {
  return (
    <section className="panel">
      <header className="panel-header">
        <ArrowRightLeft size={18} />
        <h2>Quote and Transaction</h2>
      </header>

      <div className="field-grid">
        <label className="form-group">
          <span>Account</span>
          <input
            className="input"
            min="1"
            type="number"
            value={accountID}
            onChange={(event) => onAccountIDChange(Number(event.target.value))}
          />
        </label>
        <label className="form-group">
          <span>From</span>
          <select
            className="select"
            value={inputToken}
            onChange={(event) => onInputTokenChange(event.target.value)}
          >
            {tokens.map((token) => (
              <option key={token} value={token}>
                {token}
              </option>
            ))}
          </select>
        </label>
        <label className="form-group">
          <span>To</span>
          <select
            className="select"
            value={outputToken}
            onChange={(event) => onOutputTokenChange(event.target.value)}
          >
            {tokens.map((token) => (
              <option key={token} value={token}>
                {token}
              </option>
            ))}
          </select>
        </label>
        <label className="form-group">
          <span>Amount</span>
          <input
            className="input"
            min="0.000001"
            step="0.000001"
            type="number"
            value={amountIn}
            onChange={(event) => onAmountInChange(Number(event.target.value))}
          />
        </label>
      </div>

      <div className="quote-output">
        <span>Estimated output</span>
        <strong>{estimatedOutput === null ? "Not quoted yet" : estimatedOutput.toFixed(8)}</strong>
      </div>

      <div className="actions-row">
        <button className="primary-button" type="button" onClick={onQuote}>
          <Activity size={16} />
          <span>Quote</span>
        </button>
        <button className="primary-button" type="button" onClick={onSubmit}>
          <Send size={16} />
          <span>Submit Transaction</span>
        </button>
      </div>
    </section>
  );
}
