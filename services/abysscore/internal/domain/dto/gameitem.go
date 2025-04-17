package dto

import "abysscore/internal/infrastructure/ent"

type GameItemDTO struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Collection string `json:"collection"`
	Type       int    `json:"type"`
	Rarity     int    `json:"rarity"`
}

func MapToGameItemDTOFromEnt(g *ent.GameItem) *GameItemDTO {
	if g == nil {
		return nil
	}

	return &GameItemDTO{
		ID:         g.ID,
		Name:       g.Name,
		Collection: g.Collection,
		Type:       g.Type,
		Rarity:     g.Rarity,
	}
}
