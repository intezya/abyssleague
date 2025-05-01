package mailmessage

import (
	"fmt"
	"github.com/intezya/pkglib/generate"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type LinkEmailCodeData struct {
	Message *Message `json:"-"`

	// Stored in cache
	UserID           int       `json:"user_id"`
	VerificationCode string    `json:"verification_code"`
	EmailForLink     string    `json:"email_for_link"`
	CreatedAt        time.Time `json:"created_at"`
}

// MarshalBinary implements encoding.BinaryMarshaler
func (d *LinkEmailCodeData) MarshalBinary() ([]byte, error) {
	return jsoniter.Marshal(d)
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (d *LinkEmailCodeData) UnmarshalBinary(data []byte) error {
	return jsoniter.Unmarshal(data, d)
}

func NewLinkEmailCodeMail(
	userID int,
	emailForLink string,
	validMinutes int,
) *LinkEmailCodeData {
	const verificationCodeCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789"

	verificationCode := generate.RandomString(6, verificationCodeCharset)

	const subject = "Email address confirmation"
	const mime = "text/html; charset=UTF-8"

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
		CreatedAt:        time.Now(),
	}
}
