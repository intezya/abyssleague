package examples

type InventoryItemNotFoundResponse struct {
	Message string `json:"message" example:"item not found"`
	Detail  string `json:"detail"`
	Code    int    `json:"code"    example:"404"`
	Path    string `json:"path"`
}
