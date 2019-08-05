package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	UserID       int64
	Username     string
	PasswordHash string
}

func (user *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return false
	}
	return true
}
