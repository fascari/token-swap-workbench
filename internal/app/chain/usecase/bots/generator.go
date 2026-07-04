package bots

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
)

const (
	botInterval         = 750 * time.Millisecond
	defaultAccountCount = 10
	maxSendAmount       = 25
	maxSwapAmount       = 15
)

func newBotTicker() *time.Ticker {
	return time.NewTicker(botInterval)
}

func submitRandom(ctx context.Context, client Client, accountID uint32) (operationKind, error) {
	if rand.IntN(2) == 0 {
		return operationSend, client.SubmitTransaction(ctx, domain.TransactionSubmission{
			Kind: domain.TransactionKindSend,
			Send: randomSend(accountID),
		})
	}

	return operationSwap, client.SubmitTransaction(ctx, domain.TransactionSubmission{
		Kind: domain.TransactionKindSwap,
		Swap: randomSwap(accountID),
	})
}

func randomSend(accountID uint32) domain.Send {
	return domain.Send{
		From:   accountID,
		To:     randomRecipient(accountID),
		Token:  randomToken(),
		Amount: randomAmount(maxSendAmount),
	}
}

func randomSwap(accountID uint32) domain.Swap {
	inToken, outToken := randomTokenPair()
	return domain.Swap{
		AccountID: accountID,
		InToken:   inToken,
		OutToken:  outToken,
		AmountIn:  randomAmount(maxSwapAmount),
	}
}

func randomAccountID() uint32 {
	return uint32(rand.IntN(defaultAccountCount) + 1)
}

func randomRecipient(from uint32) uint32 {
	to := randomAccountID()
	for to == from {
		to = randomAccountID()
	}
	return to
}

func randomAmount(maxAmount int) float64 {
	return float64(rand.IntN(maxAmount*100)+1) / 100
}

func randomTokenPair() (inToken domain.Token, outToken domain.Token) {
	inToken = randomToken()
	outToken = randomToken()
	for outToken == inToken {
		outToken = randomToken()
	}
	return inToken, outToken
}

func randomToken() domain.Token {
	tokens := []domain.Token{"NEX", "ETH", "USDC", "HYPE"}
	return tokens[rand.IntN(len(tokens))]
}
