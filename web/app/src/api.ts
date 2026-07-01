export type ChainStatus = {
  status: string;
};

export type QuoteResponse = {
  amount_out: number;
};

export type SwapResponse = {
  status: string;
};

export type BlockResponse = {
  id: number;
  timestamp: number;
  transactions: unknown[];
};

type SwapRequest = {
  account_id: number;
  in_token: string;
  out_token: string;
  amount_in: number;
};

async function requestJSON<T>(
  input: RequestInfo | URL,
  init?: RequestInit,
): Promise<T> {
  const response = await fetch(input, {
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {}),
    },
    ...init,
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

export function submitSwap(payload: SwapRequest): Promise<SwapResponse> {
  return requestJSON<SwapResponse>("/v1/swaps", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function fetchBlocks(count: number): Promise<BlockResponse[]> {
  return requestJSON<BlockResponse[]>(`/v1/blocks?n=${count}`);
}
