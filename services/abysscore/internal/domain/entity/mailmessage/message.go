package mailmessage

type Message struct {
	subject string
	mime    string
	body    string
}

func NewMessage(subject string, mime string, body string) *Message {
	return &Message{subject: subject, mime: mime, body: body}
}

func (m *Message) AsBytes() []byte {
	return []byte(m.subject + m.mime + m.body)
}
