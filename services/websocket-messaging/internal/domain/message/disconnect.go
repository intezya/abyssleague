package message

import "encoding/json"

type message struct {
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
	Message string `json:"message"`
}

var DisconnectByOtherClient []byte

func init() {
	var err error

	DisconnectByOtherClient, err = json.Marshal(
		message{
			Type:    "disconnect",
			Subtype: "other_client",
			Message: "You have been disconnected by another connection",
		},
	)
	if err != nil {
		panic("failed to marshal disconnect message: " + err.Error())
	}
}
