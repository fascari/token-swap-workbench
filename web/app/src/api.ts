export type ChainStatus = {
  status: string;
};

export type QuoteResponse = {
  amount_out: number;
};

export type TransactionRequest = {
  account_id: number;
  in_token: string;
  out_token: string;
  amount_in: number;
};

export type TransactionResponse = {
  status: string;
};

export type BotAction = "create" | "stop";

export type BotRequest = {
  action: BotAction;
  amount?: number;
  all?: boolean;
};

export type BotResponse = {
  status: string;
  action: string;
  requested_amount: number;
  all: boolean;
  active_bots: number;
  created_bots: number;
  stopped_bots: number;
  attempted_operations: number;
  accepted_operations: number;
  failed_operations: number;
  send_operations: number;
  swap_operations: number;
};

export type Block = {
  id: number;
  timestamp: number;
  transactions: unknown[];
};

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
  const query = new URLSearchParams({
    in: inToken,
    out: outToken,
    amount: String(amount),
  });

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
