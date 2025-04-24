package examples

type GameItemNotFound struct {
	Message string `json:"message" example:"game item not found"`
	Detail  string `json:"detail"`
	Code    int    `json:"code"    example:"404"`
	Path    string `json:"path"`
}
