package dto

// ChangePasswordDTO holds data for password change requests
type ChangePasswordDTO struct {
	Username    string
	OldPassword string
	NewPassword string
}

// CredentialsDTO holds user authentication input data
type CredentialsDTO struct {
	Username string
	Password string
	Hwid     string
}

func NewCredentialsDTO(username string, password string, hwid string) *CredentialsDTO {
	return &CredentialsDTO{Username: username, Password: password, Hwid: hwid}
}
