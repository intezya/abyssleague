package examples

import (
	"abysscore/internal/domain/dto"
	domainservice "abysscore/internal/domain/service"
)

/* TEMPLATE
type SuccessResponse struct {
	Message string `json:"message" example:"success"`
	Data    T      `json:"data"`
	Code    int    `json:"code" example:"200"`
	Path    string `json:"path"`
}
*/

type CreateGameItemDTOSuccessResponse struct {
	Message string          `json:"message" example:"success"`
	Data    dto.GameItemDTO `json:"data"`
	Code    int             `json:"code" example:"200"`
	Path    string          `json:"path"`
}

type FindGameItemDTOSuccessResponse struct {
	Message string          `json:"message" example:"success"`
	Data    dto.GameItemDTO `json:"data"`
	Code    int             `json:"code" example:"200"`
	Path    string          `json:"path"`
}

type AuthenticationSuccessResponse struct {
	Message string                             `json:"message" example:"success"`
	Data    domainservice.AuthenticationResult `json:"data"`
	Code    int                                `json:"code" example:"200"`
	Path    string                             `json:"path"`
}
