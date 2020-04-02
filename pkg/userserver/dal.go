package userserver

// UserDAL is the Data Abstraction Layer for the User Server
type UserDAL interface {
	getUserByEmail(email string) *User
	getUserById(id int) *User
	InsertUser(user User) (*User, error)
}
