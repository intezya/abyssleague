package message

import "encoding/json"

type message struct {
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
	Message string `json:"message"`
}

var DisconnectByOtherClient, _ = json.Marshal(
	message{
		Type:    "disconnect",
		Subtype: "other_client",
		Message: "You have been disconnected by another connection",
	},
)
