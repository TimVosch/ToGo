package todoserver

// TodoRepository is an interface for retrieving data
type TodoRepository interface {
	GetTodoByID(id int64) *TodoEntry
	GetTodosForUser(userID int64) []TodoEntry
	InsertTodo(todo TodoEntry) (*TodoEntry, error)
	DeleteTodo(id int64) error
}
