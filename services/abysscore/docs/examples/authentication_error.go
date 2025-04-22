package examples

type UsernameConflictResponse struct {
	Message string `json:"message" example:"user already exists"`
	Detail  string `json:"detail"`
	Code    int    `json:"code" example:"409"`
	Path    string `json:"path"`
}

type HardwareIDConflictResponse struct {
	Message string `json:"message" example:"user hwid conflict"`
	Detail  string `json:"detail"`
	Code    int    `json:"code" example:"409"`
	Path    string `json:"path"`
}

type UserNotFoundResponse struct {
	Message string `json:"message" example:"user not found"`
	Detail  string `json:"detail"`
	Code    int    `json:"code" example:"404"`
	Path    string `json:"path"`
}

type UserWrongPasswordResponse struct {
	Message string `json:"message" example:"wrong password"`
	Detail  string `json:"detail"`
	Code    int    `json:"code" example:"401"`
	Path    string `json:"path"`
}

type UserWrongHardwareIDResponse struct {
	Message string `json:"message" example:"wrong hwid"`
	Detail  string `json:"detail"`
	Code    int    `json:"code" example:"401"`
	Path    string `json:"path"`
}
