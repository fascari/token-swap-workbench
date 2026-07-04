package bots

import (
	"context"
	"sync"
)

const (
	operationSend operationKind = "send"
	operationSwap operationKind = "swap"
)

type (
	bot struct {
		cancel context.CancelFunc
	}

	manager struct {
		mu     sync.Mutex
		nextID int
		bots   map[int]bot
		stats  stats
	}

	stats struct {
		AttemptedOperations int
		AcceptedOperations  int
		FailedOperations    int
		SendOperations      int
		SwapOperations      int
	}

	operationKind string
)

func newManager() *manager {
	return &manager{bots: make(map[int]bot)}
}

func (m *manager) create(ctx context.Context, client Client, amount int) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	for range amount {
		botCtx, cancel := context.WithCancel(context.WithoutCancel(ctx))
		id := m.nextID
		m.nextID++
		m.bots[id] = bot{cancel: cancel}
		go runBot(botCtx, client, m)
	}

	return amount
}

func (m *manager) stop(amount int) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	stopped := 0
	for id, bot := range m.bots {
		if stopped == amount {
			break
		}
		bot.cancel()
		delete(m.bots, id)
		stopped++
	}

	return stopped
}

func (m *manager) stopAll() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	stopped := len(m.bots)
	for id, bot := range m.bots {
		bot.cancel()
		delete(m.bots, id)
	}

	return stopped
}

func (m *manager) record(kind operationKind, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stats.AttemptedOperations++
	if err != nil {
		m.stats.FailedOperations++
		return
	}

	m.stats.AcceptedOperations++
	switch kind {
	case operationSend:
		m.stats.SendOperations++
	case operationSwap:
		m.stats.SwapOperations++
	default:
		m.stats.FailedOperations++
	}
}

func (m *manager) output(input Input, createdBots int, stoppedBots int) Output {
	m.mu.Lock()
	defer m.mu.Unlock()

	status := statusStopped
	if len(m.bots) > 0 {
		status = statusRunning
	}

	return Output{
		Status:              status,
		Action:              input.Action,
		RequestedAmount:     input.Amount,
		All:                 input.All,
		ActiveBots:          len(m.bots),
		CreatedBots:         createdBots,
		StoppedBots:         stoppedBots,
		AttemptedOperations: m.stats.AttemptedOperations,
		AcceptedOperations:  m.stats.AcceptedOperations,
		FailedOperations:    m.stats.FailedOperations,
		SendOperations:      m.stats.SendOperations,
		SwapOperations:      m.stats.SwapOperations,
	}
}

func runBot(ctx context.Context, client Client, manager *manager) {
	ticker := newBotTicker()
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			kind, err := submitRandom(ctx, client, randomAccountID())
			manager.record(kind, err)
		}
	}
}
