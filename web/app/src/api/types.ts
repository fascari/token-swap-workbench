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
export type BotMode = "create" | "stop";

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
