package modules

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
	blockhandler "github.com/fascari/token-swap-workbench/internal/app/chain/handler/listblocks"
	quotehandler "github.com/fascari/token-swap-workbench/internal/app/chain/handler/quote"
	statushandler "github.com/fascari/token-swap-workbench/internal/app/chain/handler/status"
	swaphandler "github.com/fascari/token-swap-workbench/internal/app/chain/handler/submitswap"
	blockusecase "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/listblocks"
	quoteusecase "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/quote"
	statususecase "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/status"
	swapusecase "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/submitswap"
	"github.com/fascari/token-swap-workbench/internal/chainclient"
)

type (
	ChainModule struct {
		statusHandler statushandler.Handler
		quoteHandler  quotehandler.Handler
		swapHandler   swaphandler.Handler
		blockHandler  blockhandler.Handler
	}

	chainAdapter struct {
		client *chainclient.Client
	}
)

func NewChainModule(client *chainclient.Client) *ChainModule {
	adapter := chainAdapter{client: client}

	return &ChainModule{
		statusHandler: statushandler.New(statususecase.NewUseCase(adapter)),
		quoteHandler:  quotehandler.New(quoteusecase.NewUseCase(adapter)),
		swapHandler:   swaphandler.New(swapusecase.NewUseCase(adapter)),
		blockHandler:  blockhandler.New(blockusecase.NewUseCase(adapter)),
	}
}

func (m *ChainModule) Register(r chi.Router) {
	statushandler.RegisterRoutes(r, m.statusHandler)
	quotehandler.RegisterRoutes(r, m.quoteHandler)
	swaphandler.RegisterRoutes(r, m.swapHandler)
	blockhandler.RegisterRoutes(r, m.blockHandler)
}

func (a chainAdapter) Status(ctx context.Context) error {
	if err := a.client.Status(ctx); err != nil {
		return fmt.Errorf("requesting chain status: %w", err)
	}

	return nil
}

func (a chainAdapter) Quote(ctx context.Context, req domain.QuoteRequest) (domain.Quote, error) {
	quote, err := a.client.Quote(ctx, chainclient.QuoteRequest{
		InToken:  chainclient.Token(req.InToken),
		OutToken: chainclient.Token(req.OutToken),
		Amount:   req.Amount,
	})
	if err != nil {
		return domain.Quote{}, fmt.Errorf("requesting chain quote: %w", err)
	}

	return domain.Quote{AmountOut: quote.AmountOut}, nil
}

func (a chainAdapter) SubmitSwap(ctx context.Context, swap domain.Swap) error {
	err := a.client.SubmitSwap(ctx, chainclient.SwapRequest{
		AccountID: swap.AccountID,
		InToken:   chainclient.Token(swap.InToken),
		OutToken:  chainclient.Token(swap.OutToken),
		AmountIn:  swap.AmountIn,
	})
	if err != nil {
		return fmt.Errorf("submitting chain swap: %w", err)
	}

	return nil
}

func (a chainAdapter) Blocks(ctx context.Context, n int) ([]domain.Block, error) {
	blocks, err := a.client.Blocks(ctx, n)
	if err != nil {
		return nil, fmt.Errorf("requesting chain blocks: %w", err)
	}

	return toDomainBlocks(blocks), nil
}

func toDomainBlocks(blocks []chainclient.Block) []domain.Block {
	result := make([]domain.Block, 0, len(blocks))
	for _, block := range blocks {
		result = append(result, toDomainBlock(block))
	}

	return result
}

func toDomainBlock(block chainclient.Block) domain.Block {
	transactions := make([]domain.Transaction, 0, len(block.Transactions))
	for _, transaction := range block.Transactions {
		transactions = append(transactions, toDomainTransaction(transaction))
	}

	return domain.Block{
		ID:           block.ID,
		Timestamp:    block.Timestamp,
		Transactions: transactions,
	}
}

func toDomainTransaction(transaction chainclient.Transaction) domain.Transaction {
	return domain.Transaction{
		Swap: toDomainSwapTransaction(transaction.Swap),
		Send: toDomainSendTransaction(transaction.Send),
	}
}

func toDomainSwapTransaction(transaction *chainclient.SwapTransaction) *domain.SwapTransaction {
	if transaction == nil {
		return nil
	}

	return &domain.SwapTransaction{
		AccountID: transaction.AccountID,
		InToken:   domain.Token(transaction.InToken),
		OutToken:  domain.Token(transaction.OutToken),
		AmountIn:  transaction.AmountIn,
	}
}

func toDomainSendTransaction(transaction *chainclient.SendTransaction) *domain.SendTransaction {
	if transaction == nil {
		return nil
	}

	return &domain.SendTransaction{
		From:   transaction.From,
		To:     transaction.To,
		Amount: transaction.Amount,
		Token:  domain.Token(transaction.Token),
	}
}
