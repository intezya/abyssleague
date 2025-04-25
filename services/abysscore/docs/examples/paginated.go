package examples

import "github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"

type PaginatedGameItemsDTOResponse struct {
	Data []dto.GameItemDTO `json:"data"`

	Page       int `json:"page"`
	Size       int `json:"size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type PaginatedInventoryItemsDTOResponse struct {
	Data []dto.InventoryItemDTO `json:"data"`

	Page       int `json:"page"`
	Size       int `json:"size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}
