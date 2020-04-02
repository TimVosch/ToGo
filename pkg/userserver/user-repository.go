package userserver

// UserRepository is the Data Abstraction Layer for the User Server
type UserRepository interface {
	getUserByEmail(email string) *User
	getUserById(id int) *User
	InsertUser(user User) (*User, error)
}
