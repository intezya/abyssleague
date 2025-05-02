package examples

type AccountAlreadyHasLinkedEmail struct {
	Message string `json:"message" example:"account already has linked email"`
	Detail  string `json:"detail"`
	Code    int    `json:"code"    example:"409"`
	Path    string `json:"path"`
}

type WrongVerificationCode struct {
	Message string `json:"message" example:"bad request"`
	Detail  string `json:"detail"  example:"wrong verification code"`
	Code    int    `json:"code"    example:"400"`
	Path    string `json:"path"`
}

type EmailConflict struct {
	Message string `json:"message" example:"someone account already has this email"`
	Detail  string `json:"detail"`
	Code    int    `json:"code"    example:"409"`
	Path    string `json:"path"`
}
