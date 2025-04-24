package entity

// TokenData contains information to be encoded in auth tokens.
type TokenData struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Hwid     string `json:"hwid"`
}
