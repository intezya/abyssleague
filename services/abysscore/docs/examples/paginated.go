package examples

import "abysscore/internal/domain/dto"

type PaginatedGameItemsDTOResponse struct {
	Data []dto.GameItemDTO `json:"data"`

	Page       int `json:"page"`
	Size       int `json:"size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}
