package todoserver

import "errors"

// TodoMemoryRepository is an in memory Database for the TodoServer
type TodoMemoryRepository struct {
	todos  []TodoEntry
	lastID int64
}

func (db *TodoMemoryRepository) nextID() int64 {
	db.lastID++
	return db.lastID
}

// NewTodoMemoryRepository instantiates a new in memory database
func NewTodoMemoryRepository() *TodoMemoryRepository {
	return &TodoMemoryRepository{
		[]TodoEntry{},
		1,
	}
}

// GetTodoByID ...
func (db *TodoMemoryRepository) GetTodoByID(id int64) *TodoEntry {
	// Try to find the todo
	for _, v := range db.todos {
		if v.ID == id {
			return &v
		}
	}
	return nil
}

// GetTodosForUser ...
func (db *TodoMemoryRepository) GetTodosForUser(id int64) []TodoEntry {
	return db.todos
}

// InsertTodo ...
func (db *TodoMemoryRepository) InsertTodo(todo TodoEntry) (*TodoEntry, error) {
	todo.ID = db.nextID()
	db.todos = append(db.todos, todo)
	return &todo, nil
}

// DeleteTodo ...
func (db *TodoMemoryRepository) DeleteTodo(id int64) error {
	for i, v := range db.todos {
		if v.ID == id {
			// Remove from slice
			db.todos = append(db.todos[:i], db.todos[i+1:]...)
			return nil
		}
	}
	return errors.New("Not found")
}
