package dto

import (
	"time"
)

type GameItemDTO struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Collection string    `json:"collection"`
	Type       int       `json:"type"`
	Rarity     int       `json:"rarity"`
	CreatedAt  time.Time `json:"created_at"`
}
