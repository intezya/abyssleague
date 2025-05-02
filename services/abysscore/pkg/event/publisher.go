package eventlib

import (
	"context"
	"fmt"
	"sync"

	"github.com/intezya/pkglib/itertools"
	"github.com/intezya/pkglib/logger"
	"go.uber.org/zap"
)

type eventBus struct {
	handlers []Handler
	events   chan ApplicationEvent
	cancel   context.CancelFunc
}

type ApplicationEventPublisher struct {
	mu          sync.RWMutex
	subscribers map[string]*eventBus
	workerCount int
	bufferSize  int
}

func NewApplicationEventPublisher(workerCount int, bufferSize int) *ApplicationEventPublisher {
	if workerCount <= 0 {
		workerCount = 1
	}

	if bufferSize <= 0 {
		bufferSize = 100
	}

	return &ApplicationEventPublisher{
		subscribers: make(map[string]*eventBus),
		workerCount: workerCount,
		bufferSize:  bufferSize,
	}
}

func (p *ApplicationEventPublisher) Register(
	event ApplicationEvent,
	handler Handler,
	middleware ...Middleware,
) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, mw := range middleware {
		handler = mw(handler)
	}

	eventType := typeName(event)
	bus, exists := p.subscribers[eventType]

	if !exists {
		ctx, cancel := context.WithCancel(context.Background())

		bus = &eventBus{
			handlers: []Handler{},
			events:   make(chan ApplicationEvent, p.bufferSize),
			cancel:   cancel,
		}

		for range p.workerCount {
			go p.startWorker(ctx, bus)
		}

		p.subscribers[eventType] = bus
	}

	bus.handlers = append(bus.handlers, handler)
}

func (p *ApplicationEventPublisher) Unregister(event ApplicationEvent) {
	p.mu.Lock()
	defer p.mu.Unlock()

	eventType := typeName(event)

	bus, exists := p.subscribers[eventType]
	if exists {
		bus.cancel()
		close(bus.events)
		delete(p.subscribers, eventType)
	}
}

func (p *ApplicationEventPublisher) Publish(event ApplicationEvent) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	eventType := typeName(event)

	bus, exists := p.subscribers[eventType]

	if exists {
		select {
		case bus.events <- event:
		default:
			logger.Log.Infoln("event bus is full")
			// TODO: republish
		}
	} else {
		availableBuses := itertools.GetMapKeys(p.subscribers)

		logger.Log.Warnln("Event bus not found for event", eventType, zap.Any("available_buses", availableBuses))
	}
}

func (p *ApplicationEventPublisher) startWorker(ctx context.Context, bus *eventBus) {
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-bus.events:
			if !ok {
				return
			}

			for _, handler := range bus.handlers {
				handler(event)
			}
		}
	}
}

func typeName(event ApplicationEvent) string {
	return fmt.Sprintf("%T", event)
}
