package todoserver

// TodoRepository is an interface for retrieving data
type TodoRepository interface {
	GetTodoByID(id int) *TodoEntry
	GetTodosForUser(userID int) []TodoEntry
	InsertTodo(todo TodoEntry) (*TodoEntry, error)
	DeleteTodo(id int) error
}
