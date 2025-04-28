package mailmessage

import "github.com/intezya/pkglib/generate"

type LinkEmailCodeData struct {
	UserID           int    `json:"user_id"`
	EmailForLink     string `json:"email_for_link"`
	VerificationCode string `json:"verification_code"`
}

func NewLinkEmailCodeMail(
	userID int,
	emailForLink string,
) *LinkEmailCodeData {
	const verificationCodeCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789"

	return &LinkEmailCodeData{
		UserID:           userID,
		EmailForLink:     emailForLink,
		VerificationCode: generate.RandomString(6, verificationCodeCharset),
	}
}
