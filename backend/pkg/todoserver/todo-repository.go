package todoserver

// TodoRepository is an interface for retrieving data
type TodoRepository interface {
	GetTodoByID(id interface{}) *TodoEntry
	GetTodosForUser(userID interface{}) []TodoEntry
	InsertTodo(todo TodoEntry) (*TodoEntry, error)
	DeleteTodo(id interface{}) error
}
