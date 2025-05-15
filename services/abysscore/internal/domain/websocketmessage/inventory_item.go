package websocketmessage

import "github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"

const (
	inventoryMessageType = "inventory"
	itemMessageSubtype   = "item"
)

type InventoryItemObtainedMessage struct {
	*BaseMessage

	Data struct {
		Item *dto.InventoryItemDTO `json:"item"`
	} `json:"data"`
}

func NewInventoryItemObtainedMessage(
	eventID string,
	performerName string,
	item *dto.InventoryItemDTO,
) *InventoryItemObtainedMessage {
	const message = "item obtained"

	return &InventoryItemObtainedMessage{
		BaseMessage: NewBaseMessage(
			eventID,
			inventoryMessageType,
			itemMessageSubtype,
			message,
			performerName,
		),
		Data: struct {
			Item *dto.InventoryItemDTO `json:"item"`
		}{
			Item: item,
		},
	}
}

type InventoryItemRevokedMessage struct {
	*BaseMessage

	Data struct {
		Item *dto.InventoryItemDTO `json:"item"`
	} `json:"data"`
}

func NewInventoryItemRevokedMessage(
	eventID string,
	performerName string,
	item *dto.InventoryItemDTO,
) *InventoryItemRevokedMessage {
	const message = "item revoked"

	return &InventoryItemRevokedMessage{
		BaseMessage: NewBaseMessage(
			eventID,
			inventoryMessageType,
			itemMessageSubtype,
			message,
			performerName,
		),
		Data: struct {
			Item *dto.InventoryItemDTO `json:"item"`
		}{
			Item: item,
		},
	}
}
