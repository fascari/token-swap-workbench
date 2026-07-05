package chainclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
)

func (c Client) SubmitTransaction(ctx context.Context, transaction domain.TransactionSubmission) error {
	payload, err := transactionPayload(transaction)
	if err != nil {
		return fmt.Errorf("preparing chain transaction: %w", err)
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(payload); err != nil {
		return classify(fmt.Errorf("encoding chain transaction: %w", err))
	}

	endpoint := c.endpoint("transaction")
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), &body)
	if err != nil {
		return classify(fmt.Errorf("building chain transaction request: %w", err))
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return classify(fmt.Errorf("submitting chain transaction: %w", err))
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if err := ensureSuccess(resp); err != nil {
		return classify(fmt.Errorf("submitting chain transaction: %w", err))
	}

	return nil
}

func transactionPayload(transaction domain.TransactionSubmission) (transactionEnvelope, error) {
	switch transaction.Kind {
	case domain.TransactionKindSend:
		return transactionEnvelope{Send: sendPayload(transaction.Send)}, nil
	case domain.TransactionKindSwap:
		return transactionEnvelope{Swap: swapPayload(transaction.Swap)}, nil
	default:
		return transactionEnvelope{}, fmt.Errorf("unsupported transaction kind %q", transaction.Kind)
	}
}

func sendPayload(send domain.Send) *sendRequest {
	return new(sendRequest{
		From:   send.From,
		To:     send.To,
		Token:  send.Token,
		Amount: send.Amount,
	})
}

func swapPayload(swap domain.Swap) *swapRequest {
	return new(swapRequest{
		AccountID: swap.AccountID,
		InToken:   swap.InToken,
		OutToken:  swap.OutToken,
		AmountIn:  swap.AmountIn,
	})
}
