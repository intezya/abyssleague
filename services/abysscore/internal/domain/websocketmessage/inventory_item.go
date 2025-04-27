package websocketmessage

import "github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"

const inventoryMessageType = "inventory"

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
	const subtype = "new_item"
	const message = "item obtained"

	return &InventoryItemObtainedMessage{
		BaseMessage: NewBaseMessage(
			eventID,
			inventoryMessageType,
			subtype,
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
