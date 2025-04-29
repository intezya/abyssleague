package mailmessage

import (
	"fmt"
	"github.com/intezya/pkglib/generate"
	"time"
)

type LinkEmailCodeData struct {
	*Message

	// Stored in cache
	UserID           int    `json:"user_id"`
	VerificationCode string `json:"verification_code"`
	EmailForLink     string `json:"email_for_link"`
}

func NewLinkEmailCodeMail(
	userID int,
	emailForLink string,
	validMinutes int,
) *LinkEmailCodeData {
	const verificationCodeCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789"

	verificationCode := generate.RandomString(6, verificationCodeCharset)

	const subject = "Email address confirmation"
	const mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(
		linkEmailCodeMessageBodyTemplate,
		emailForLink,
		verificationCode,
		validMinutes,
		time.Now().Year(),
	)

	return &LinkEmailCodeData{
		UserID:           userID,
		VerificationCode: verificationCode,
		EmailForLink:     emailForLink,
		Message:          NewMessage(subject, mime, body),
	}
}
