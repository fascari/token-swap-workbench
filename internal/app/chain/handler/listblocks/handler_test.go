package listblocks_test

import (
	"errors"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/fascari/token-swap-workbench/internal/app/chain/handler/listblocks"
	listblocksfixtures "github.com/fascari/token-swap-workbench/internal/app/chain/handler/listblocks/testdata"
	chaintestdata "github.com/fascari/token-swap-workbench/internal/app/chain/testdata"
	listblocksuc "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/listblocks"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/listblocks/mocks"
	"github.com/fascari/token-swap-workbench/pkg/handlertest"
)

const (
	blockCount = 2
)

type (
	ListBlocksSuite struct {
		handlertest.Suite
	}
)

func TestListBlocksSuite(t *testing.T) {
	suite.Run(t, new(ListBlocksSuite))
}

func (s *ListBlocksSuite) TestHandler_Handle_ShouldReturnBlocks() {
	client := mocks.NewClient(s.T())
	client.EXPECT().Blocks(s.T().Context(), blockCount).Return(chaintestdata.Blocks(), nil)
	handler := listblocks.New(listblocksuc.New(client))

	s.Serve(handler.Handle, http.MethodGet, "/chain/blocks", handlertest.WithQuery("n", strconv.Itoa(blockCount)))

	s.RequireJSONResponse(http.StatusOK, listblocksfixtures.BlocksResponse)
}

func (s *ListBlocksSuite) TestHandler_Handle_ShouldReturnBadRequestWhenCountIsInvalid() {
	handler := listblocks.New(listblocksuc.New(mocks.NewClient(s.T())))

	s.Serve(handler.Handle, http.MethodGet, "/chain/blocks", handlertest.WithQuery("n", "bad"))

	s.RequireStatus(http.StatusBadRequest)
}

func (s *ListBlocksSuite) TestHandler_Handle_ShouldReturnBadGatewayWhenClientListFails() {
	client := mocks.NewClient(s.T())
	client.EXPECT().Blocks(s.T().Context(), blockCount).Return(nil, errors.New("blocks unavailable"))
	handler := listblocks.New(listblocksuc.New(client))

	s.Serve(handler.Handle, http.MethodGet, "/chain/blocks", handlertest.WithQuery("n", strconv.Itoa(blockCount)))

	s.RequireStatus(http.StatusBadGateway)
}
