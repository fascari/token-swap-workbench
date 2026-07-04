package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
)

type Client struct {
	mock.Mock
}

type Client_Expecter struct {
	mock *mock.Mock
}

type Client_SubmitTransaction_Call struct {
	*mock.Call
}

func NewClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *Client {
	m := &Client{}
	m.Mock.Test(t)

	t.Cleanup(func() { m.AssertExpectations(t) })

	return m
}

func (m *Client) EXPECT() *Client_Expecter {
	return &Client_Expecter{mock: &m.Mock}
}

func (m *Client) SubmitTransaction(ctx context.Context, transaction domain.TransactionSubmission) error {
	ret := m.Called(ctx, transaction)

	if fn, ok := ret.Get(0).(func(context.Context, domain.TransactionSubmission) error); ok {
		return fn(ctx, transaction)
	}

	return ret.Error(0)
}

func (e *Client_Expecter) SubmitTransaction(ctx any, transaction any) *Client_SubmitTransaction_Call {
	return &Client_SubmitTransaction_Call{Call: e.mock.On("SubmitTransaction", ctx, transaction)}
}

func (c *Client_SubmitTransaction_Call) Run(run func(context.Context, domain.TransactionSubmission)) *Client_SubmitTransaction_Call {
	c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.TransactionSubmission))
	})
	return c
}

func (c *Client_SubmitTransaction_Call) Return(err error) *Client_SubmitTransaction_Call {
	c.Call.Return(err)
	return c
}

func (c *Client_SubmitTransaction_Call) RunAndReturn(run func(context.Context, domain.TransactionSubmission) error) *Client_SubmitTransaction_Call {
	c.Call.Return(run)
	return c
}
