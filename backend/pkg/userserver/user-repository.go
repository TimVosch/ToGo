package userserver

// UserRepository is the Data Abstraction Layer for the User Server
type UserRepository interface {
	GetUserByEmail(email string) *User
	GetUserByID(id interface{}) *User
	InsertUser(user User) (*User, error)
}
