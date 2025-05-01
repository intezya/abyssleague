package eventlib

const (
	EveryoneIsReceiver = 0
	SystemIsSender     = 0
)

type ApplicationEvent interface {
	EventID() string
	ReceiverID() int
	SenderID() int
}

type (
	Handler    func(event ApplicationEvent)
	Middleware func(next Handler) Handler
)

type Publisher interface {
	Publish(event ApplicationEvent)
}

type Manager interface {
	Register(event ApplicationEvent, handler Handler, middleware ...Middleware)
	Unregister(event ApplicationEvent)
}
