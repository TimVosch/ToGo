package userserver

import "golang.org/x/crypto/bcrypt"

// User represent a system account
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SetPassword sets and hashes the given password
func (u *User) SetPassword(password string) {
	out, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		//
	}
	u.Password = string(out)
}
