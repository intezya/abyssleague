package dto

type PaginatedResult[T any] struct {
	Data []T `json:"data"`

	Page       int `json:"page"`
	Size       int `json:"size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}
