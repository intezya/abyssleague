package request

import (
	"abysscore/internal/domain/dto"
	"time"
)

type CreateUpdateGameItem struct {
	Name       string `json:"name"       validate:"required"`
	Collection string `json:"collection" validate:"required"`
	Type       int    `json:"type"       validate:"required"`
	Rarity     int    `json:"rarity"     validate:"required"`
}

func (c *CreateUpdateGameItem) ToDTO() *dto.GameItemDTO { // TODO: repository should accept CreateUpdateGameItemDTO from domain
	return &dto.GameItemDTO{
		ID:         0, // blank
		Name:       c.Name,
		Collection: c.Collection,
		Type:       c.Type,
		Rarity:     c.Rarity,
		CreatedAt:  time.Time{}, // blank
	}
}
