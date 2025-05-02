package examples

type UsernameConflictResponse struct {
	Message string `json:"message" example:"user already exists"`
	Detail  string `json:"detail"`
	Code    int    `json:"code"    example:"409"`
	Path    string `json:"path"`
}

type HardwareIDConflictResponse struct {
	Message string `json:"message" example:"user hardware id conflict"`
	Detail  string `json:"detail"`
	Code    int    `json:"code"    example:"409"`
	Path    string `json:"path"`
}

type UserNotFoundResponse struct {
	Message string `json:"message" example:"user not found"`
	Detail  string `json:"detail"`
	Code    int    `json:"code"    example:"404"`
	Path    string `json:"path"`
}

type UserWrongPasswordResponse struct {
	Message string `json:"message" example:"unauthorized"`
	Detail  string `json:"detail" example:"wrong password"`
	Code    int    `json:"code"    example:"401"`
	Path    string `json:"path"`
}

type UserWrongHardwareIDResponse struct {
	Message string `json:"message" example:"unauthorized"`
	Detail  string `json:"detail" example:"wrong hardware id"`
	Code    int    `json:"code"    example:"401"`
	Path    string `json:"path"`
}
