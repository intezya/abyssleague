package examples

import "github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"

type PaginatedGameItemsDTOResponse struct {
	Data []dto.GameItemDTO `json:"data"`

	Page       int `json:"page"        example:"1"`
	Size       int `json:"size"        example:"10"`
	TotalItems int `json:"total_items" example:"777"`
	TotalPages int `json:"total_pages" example:"78"`
}

type PaginatedInventoryItemsDTOResponse struct {
	Data []dto.InventoryItemDTO `json:"data"`

	Page       int `json:"page"        example:"1"`
	Size       int `json:"size"        example:"10"`
	TotalItems int `json:"total_items" example:"777"`
	TotalPages int `json:"total_pages" example:"78"`
}
