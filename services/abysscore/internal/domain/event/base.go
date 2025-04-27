package event

import (
	"github.com/google/uuid"
	"github.com/intezya/abyssleague/services/abysscore/pkg/optional"
)

type BaseEvent struct {
	id         string
	receiverID int
	senderID   int
}

func newBaseEvent(eventID optional.String, receiverID int, senderID int) *BaseEvent {
	return &BaseEvent{
		id:         eventID.DefaultFn(uuid.NewString),
		receiverID: receiverID,
		senderID:   senderID,
	}
}

func (e *BaseEvent) EventID() string {
	return e.id
}

func (e *BaseEvent) ReceiverID() int {
	return e.receiverID
}

func (e *BaseEvent) SenderID() int {
	return e.senderID
}
