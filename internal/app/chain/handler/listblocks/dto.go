package listblocks

import "github.com/fascari/token-swap-workbench/internal/app/chain/domain"

type (
	responseDTO []blockDTO

	blockDTO struct {
		ID           uint64           `json:"id"`
		Timestamp    uint64           `json:"timestamp"`
		Transactions []transactionDTO `json:"transactions"`
	}

	transactionDTO struct {
		Swap *swapTransactionDTO `json:"Swap,omitzero"`
		Send *sendTransactionDTO `json:"Send,omitzero"`
	}

	swapTransactionDTO struct {
		AccountID uint32  `json:"account_id"`
		InToken   string  `json:"in_token"`
		OutToken  string  `json:"out_token"`
		AmountIn  float64 `json:"amount_in"`
	}

	sendTransactionDTO struct {
		From   uint32  `json:"from"`
		To     uint32  `json:"to"`
		Amount float64 `json:"amount"`
		Token  string  `json:"token"`
	}
)

func toResponse(blocks []domain.Block) responseDTO {
	response := make(responseDTO, 0, len(blocks))
	for _, block := range blocks {
		response = append(response, toBlockDTO(block))
	}

	return response
}

func toBlockDTO(block domain.Block) blockDTO {
	transactions := make([]transactionDTO, 0, len(block.Transactions))
	for _, transaction := range block.Transactions {
		transactions = append(transactions, toTransactionDTO(transaction))
	}

	return blockDTO{
		ID:           block.ID,
		Timestamp:    block.Timestamp,
		Transactions: transactions,
	}
}

func toTransactionDTO(transaction domain.Transaction) transactionDTO {
	return transactionDTO{
		Swap: toSwapTransactionDTO(transaction.Swap),
		Send: toSendTransactionDTO(transaction.Send),
	}
}

func toSwapTransactionDTO(transaction *domain.SwapTransaction) *swapTransactionDTO {
	if transaction == nil {
		return nil
	}

	return new(swapTransactionDTO{
		AccountID: transaction.AccountID,
		InToken:   string(transaction.InToken),
		OutToken:  string(transaction.OutToken),
		AmountIn:  transaction.AmountIn,
	})
}

func toSendTransactionDTO(transaction *domain.SendTransaction) *sendTransactionDTO {
	if transaction == nil {
		return nil
	}

	return new(sendTransactionDTO{
		From:   transaction.From,
		To:     transaction.To,
		Amount: transaction.Amount,
		Token:  string(transaction.Token),
	})
}
