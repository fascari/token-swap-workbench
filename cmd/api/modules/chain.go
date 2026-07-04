package modules

import (
	"github.com/go-chi/chi/v5"

	bothandler "github.com/fascari/token-swap-workbench/internal/app/chain/handler/bots"
	blockhandler "github.com/fascari/token-swap-workbench/internal/app/chain/handler/listblocks"
	quotehandler "github.com/fascari/token-swap-workbench/internal/app/chain/handler/quote"
	statushandler "github.com/fascari/token-swap-workbench/internal/app/chain/handler/status"
	transactionhandler "github.com/fascari/token-swap-workbench/internal/app/chain/handler/transaction"
	botsusecase "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/bots"
	blockusecase "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/listblocks"
	quoteusecase "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/quote"
	statususecase "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/status"
	transactionusecase "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/transaction"
	"github.com/fascari/token-swap-workbench/internal/chainclient"
)

type (
	ChainModule struct {
		statusHandler      statushandler.Handler
		quoteHandler       quotehandler.Handler
		transactionHandler transactionhandler.Handler
		botHandler         bothandler.Handler
		blockHandler       blockhandler.Handler
	}
)

func NewChainModule(client *chainclient.Client) *ChainModule {
	return new(ChainModule{
		statusHandler:      statushandler.New(statususecase.New(client)),
		quoteHandler:       quotehandler.New(quoteusecase.New(client)),
		transactionHandler: transactionhandler.New(transactionusecase.New(client)),
		botHandler:         bothandler.New(botsusecase.New(client)),
		blockHandler:       blockhandler.New(blockusecase.New(client)),
	})
}

func (m *ChainModule) Register(r chi.Router) {
	statushandler.RegisterRoutes(r, m.statusHandler)
	quotehandler.RegisterRoutes(r, m.quoteHandler)
	transactionhandler.RegisterRoutes(r, m.transactionHandler)
	bothandler.RegisterRoutes(r, m.botHandler)
	blockhandler.RegisterRoutes(r, m.blockHandler)
}
