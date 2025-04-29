package request

type LinkEmailRequest struct {
	Email string `json:"email" validate:"required" example:"intezya@gmail.com"` // TODO: validate: email
}

type EnterCodeForEmailLinkRequest struct {
	VerificationCode string `json:"verification_code" validate:"required" example:"Q2JV01"`
}
