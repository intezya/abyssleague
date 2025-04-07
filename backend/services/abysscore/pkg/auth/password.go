package auth

type PasswordEncoder struct{}

func (p PasswordEncoder) Encode(raw string) (hash string) {
	//TODO implement me
	panic("implement me")
}

func (p PasswordEncoder) Verify(raw, hash string) bool {
	//TODO implement me
	panic("implement me")
}
