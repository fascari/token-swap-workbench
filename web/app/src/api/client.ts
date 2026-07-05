import type {
  Block,
  BotRequest,
  BotResponse,
  ChainStatus,
  QuoteResponse,
  TransactionRequest,
  TransactionResponse,
} from "./types";

async function requestJSON<T>(url: RequestInfo, options?: RequestInit): Promise<T> {
  const response = await fetch(url, {
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
    ...options,
  });

  if (!response.ok) {
    throw new Error(response.statusText || "Request failed");
  }

  return (await response.json()) as T;
}

export function fetchStatus(): Promise<ChainStatus> {
  return requestJSON<ChainStatus>("/v1/chain/status");
}

export function fetchQuote(
  inToken: string,
  outToken: string,
  amount: number,
): Promise<QuoteResponse> {
  const query = new URLSearchParams({ in: inToken, out: outToken, amount: String(amount) });
  return requestJSON<QuoteResponse>(`/v1/quote?${query.toString()}`);
}

export function submitTransaction(payload: TransactionRequest): Promise<TransactionResponse> {
  return requestJSON<TransactionResponse>("/v1/transactions", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function runBots(payload: BotRequest): Promise<BotResponse> {
  return requestJSON<BotResponse>("/v1/bots", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function fetchBlocks(count: number): Promise<Block[]> {
  return requestJSON<Block[]>(`/v1/blocks?n=${count}`);
}
