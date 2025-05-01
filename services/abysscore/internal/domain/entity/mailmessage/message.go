package mailmessage

type Message struct {
	Subject string
	Mime    string
	Body    string
}

func NewMessage(subject string, mime string, body string) *Message {
	return &Message{Subject: subject, Mime: mime, Body: body}
}

func (m *Message) AsBytes() []byte {
	return []byte(m.Subject + m.Mime + m.Body)
}
