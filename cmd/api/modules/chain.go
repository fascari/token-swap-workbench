package modules

import (
	"github.com/go-chi/chi/v5"

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
)

func NewChainModule(client *chainclient.Client) *ChainModule {
	return new(ChainModule{
		statusHandler: statushandler.New(statususecase.NewUseCase(client)),
		quoteHandler:  quotehandler.New(quoteusecase.NewUseCase(client)),
		swapHandler:   swaphandler.New(swapusecase.NewUseCase(client)),
		blockHandler:  blockhandler.New(blockusecase.NewUseCase(client)),
	})
}

func (m *ChainModule) Register(r chi.Router) {
	statushandler.RegisterRoutes(r, m.statusHandler)
	quotehandler.RegisterRoutes(r, m.quoteHandler)
	swaphandler.RegisterRoutes(r, m.swapHandler)
	blockhandler.RegisterRoutes(r, m.blockHandler)
}
