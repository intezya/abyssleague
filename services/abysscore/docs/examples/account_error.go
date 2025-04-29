package examples

type AccountAlreadyHasLinkedEmail struct {
	Message string `json:"message" example:"account already has linked email"`
	Detail  string `json:"detail"`
	Code    int    `json:"code"    example:"409"`
	Path    string `json:"path"`
}

type WrongVerificationCode struct {
	Message string `json:"message" example:"wrong verification code"`
	Detail  string `json:"detail"`
	Code    int    `json:"code"    example:"400"`
	Path    string `json:"path"`
}

type EmailConflict struct {
	Message string `json:"message" example:"email conflict"`
	Detail  string `json:"detail"`
	Code    int    `json:"code"    example:"409"`
	Path    string `json:"path"`
}
