package auth

type AuthenticationData interface {
	GetUsername() (username string)
	GetHardwareID() (hwid string)
}

type TokenProvider struct {
	secretKey []byte
	issuer    string
}

func NewTokenProvider(secretKey []byte) *TokenProvider {
	return &TokenProvider{secretKey: secretKey}
}

func (t TokenProvider) GenerateToken(authentication AuthenticationData) map[string]string {
	//TODO implement me
	panic("implement me")
}
