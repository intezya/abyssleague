package user

type TokenProvider interface {
	GenerateToken(user *User) map[string]string
}

type PasswordEncoder interface {
	Encode(raw string) (hash string)
	Verify(raw, hash string) bool
}
