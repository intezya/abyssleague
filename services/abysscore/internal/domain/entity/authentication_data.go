package entity

type comparator func(raw, hash string) bool

type CredentialsDTO struct {
	Username string
	Password string
	Hwid     string
}

func NewCredentialsDTO(username string, password string, hwid string) *CredentialsDTO {
	return &CredentialsDTO{Username: username, Password: password, Hwid: hwid}
}

type AuthenticationData struct {
	id       int
	username string
	password string
	hwid     *string
}

type TokenData struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Hwid     string `json:"hwid"`
}

func NewAuthenticationData(id int, username string, password string, hwid *string) *AuthenticationData {
	return &AuthenticationData{id: id, username: username, password: password, hwid: hwid}
}

func (a *AuthenticationData) ComparePassword(
	password string,
	comparator comparator,
) bool {
	return comparator(password, a.password)
}

func (a *AuthenticationData) CompareHWID(
	hwid string,
	comparator comparator,
) (ok bool, needsUpdate bool) {
	if a.hwid == nil {
		return true, true
	}

	return comparator(hwid, *a.hwid), false
}

func (a *AuthenticationData) TokenData() *TokenData {
	hwid := ""

	if a.hwid != nil {
		hwid = *a.hwid
	}

	return &TokenData{
		ID:       a.id,
		Username: a.username,
		Hwid:     hwid,
	}
}

func (a *AuthenticationData) SetHWID(hwid string) {
	a.hwid = &hwid
}

func (a *AuthenticationData) UserID() int {
	return a.id
}

type ChangePasswordDTO struct {
	Username    string
	OldPassword string
	NewPassword string
}
