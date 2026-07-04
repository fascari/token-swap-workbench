package chainclient

import (
	"fmt"
	"net/http"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
)

type (
	quoteResponse struct {
		AmountOut float64 `json:"amount_out"`
	}

	swapRequest struct {
		AccountID uint32       `json:"account"`
		InToken   domain.Token `json:"in_token"`
		OutToken  domain.Token `json:"out_token"`
		AmountIn  float64      `json:"amount_in"`
	}
	sendRequest struct {
		From   uint32       `json:"from"`
		To     uint32       `json:"to"`
		Amount float64      `json:"amount"`
		Token  domain.Token `json:"token"`
	}
	block struct {
		ID           uint64        `json:"id"`
		Timestamp    uint64        `json:"timestamp"`
		Transactions []transaction `json:"transactions"`
	}

	transaction struct {
		Swap *swapTransaction `json:"Swap,omitzero"`
		Send *sendTransaction `json:"Send,omitzero"`
	}

	swapTransaction struct {
		AccountID uint32       `json:"account_id"`
		InToken   domain.Token `json:"in_token"`
		OutToken  domain.Token `json:"out_token"`
		AmountIn  float64      `json:"amount_in"`
	}

	sendTransaction struct {
		From   uint32       `json:"from"`
		To     uint32       `json:"to"`
		Amount float64      `json:"amount"`
		Token  domain.Token `json:"token"`
	}

	responseError struct {
		statusCode int
		status     string
		detail     string
	}

	transactionEnvelope struct {
		Send *sendRequest `json:"Send,omitempty"`
		Swap *swapRequest `json:"Swap,omitempty"`
	}
)

func (e *responseError) Error() string {
	if e.detail == "" {
		return fmt.Sprintf("chain returned %s", e.status)
	}

	return fmt.Sprintf("chain returned %s: %s", e.status, e.detail)
}

func (e *responseError) clientRejected() bool {
	return e.statusCode >= http.StatusBadRequest && e.statusCode < http.StatusInternalServerError
}

func toDomainBlocks(blocks []block) []domain.Block {
	result := make([]domain.Block, 0, len(blocks))
	for _, b := range blocks {
		result = append(result, toDomainBlock(b))
	}

	return result
}

func toDomainBlock(b block) domain.Block {
	transactions := make([]domain.Transaction, 0, len(b.Transactions))
	for _, tx := range b.Transactions {
		transactions = append(transactions, toDomainTransaction(tx))
	}

	return domain.Block{
		ID:           b.ID,
		Timestamp:    b.Timestamp,
		Transactions: transactions,
	}
}

func toDomainTransaction(tx transaction) domain.Transaction {
	return domain.Transaction{
		Swap: toDomainSwapTransaction(tx.Swap),
		Send: toDomainSendTransaction(tx.Send),
	}
}

func toDomainSwapTransaction(tx *swapTransaction) *domain.SwapTransaction {
	if tx == nil {
		return nil
	}

	return new(domain.SwapTransaction{
		AccountID: tx.AccountID,
		InToken:   tx.InToken,
		OutToken:  tx.OutToken,
		AmountIn:  tx.AmountIn,
	})
}

func toDomainSendTransaction(tx *sendTransaction) *domain.SendTransaction {
	if tx == nil {
		return nil
	}

	return new(domain.SendTransaction{
		From:   tx.From,
		To:     tx.To,
		Amount: tx.Amount,
		Token:  tx.Token,
	})
}
