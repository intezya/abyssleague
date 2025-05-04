package request

import (
	"time"

	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
)

type CreateUpdateGameItem struct {
	Name       string `json:"name"       validate:"required" example:"Whirling Mark"`
	Collection string `json:"collection" validate:"required" example:"Shadow Sigils"`
	Type       int    `json:"type"       validate:"required" example:"4"`
	Rarity     int    `json:"rarity"     validate:"required" example:"4"`
}

// TODO: repository should accept CreateUpdateGameItemDTO from domain.
func (c *CreateUpdateGameItem) ToDTO() *dto.GameItemDTO {
	return &dto.GameItemDTO{
		ID:         0, // blank
		Name:       c.Name,
		Collection: c.Collection,
		Type:       c.Type,
		Rarity:     c.Rarity,
		CreatedAt:  time.Time{}, // blank
	}
}
