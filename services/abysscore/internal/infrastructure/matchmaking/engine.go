package matchmaking

import (
	"context"
	"github.com/intezya/abyssleague/services/abysscore/pkg/types"
	"time"
)

const defaultEnginePoolSize = 8
const defaultMaxPointDifferenceForPlayers = 50

type Engine interface {
	types.Runnable
	RegisterPlayer(*Entity) bool   // false if already registered
	UnregisterPlayer(*Entity) bool // false if not registered
}

// TODO: add link on tech about

type engine struct {
	pool                         []*Entity
	callback                     callback
	repeatDelay                  time.Duration
	maxPointDifferenceForPlayers uint
}

func NewEngine(callback callback, opts ...EngineOption) Engine {
	e := &engine{
		pool:                         make([]*Entity, defaultEnginePoolSize),
		callback:                     callback,
		maxPointDifferenceForPlayers: defaultMaxPointDifferenceForPlayers,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func (e *engine) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		}

		// TODO: infra logic
		// if found, use e.callback
	}
}

func (e *engine) RegisterPlayer(entity *Entity) bool {
	//TODO implement me
	panic("implement me")
}

func (e *engine) UnregisterPlayer(entity *Entity) bool {
	//TODO implement me
	panic("implement me")
}
