package todoserver

// TodoDAL is an interface for retrieving data
type TodoDAL interface {
	GetTodoByID(id int) *TodoEntry
	GetTodosForUser(userID int) []TodoEntry
	InsertTodo(todo TodoEntry) (*TodoEntry, error)
	DeleteTodo(id int) error
}
