package todoserver

// TodoDAL is an interface for retrieving data
type TodoDAL interface {
	GetTodosForUser(userID int) []TodoEntry
}
