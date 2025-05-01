package mail

import "net/smtp"

type mailDevCustomAuth struct {
	username, password string
}

func newCustomAuth(username, password string) smtp.Auth {
	return &mailDevCustomAuth{username, password}
}

func (a *mailDevCustomAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *mailDevCustomAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		}
	}
	return nil, nil
}
