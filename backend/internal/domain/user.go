package domain

type User struct {
	ID           ID
	Username     string
	PasswordHash []byte
}
