package websocketmessage

type (
	messageType    string
	messageSubtype string
)

const SystemIsSenderName = "system"

type BaseMessage struct {
	EventID string `json:"event_id"`

	Type    messageType    `json:"type"`
	Subtype messageSubtype `json:"subtype"`
	Message string         `json:"message"`

	SenderName string `json:"sender_name"` // system or username
}

func NewBaseMessage(
	eventID string,
	messageType messageType,
	messageSubtype messageSubtype,
	message string,
	senderName string,
) *BaseMessage {
	return &BaseMessage{
		EventID:    eventID,
		Type:       messageType,
		Subtype:    messageSubtype,
		Message:    message,
		SenderName: senderName,
	}
}
