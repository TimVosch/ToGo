package todoserver

// TodoRepository is an interface for retrieving data
type TodoRepository interface {
	GetTodoByID(id interface{}) *TodoEntry
	GetTodosForUser(userID int64) []TodoEntry
	InsertTodo(todo TodoEntry) (*TodoEntry, error)
	DeleteTodo(id interface{}) error
}
