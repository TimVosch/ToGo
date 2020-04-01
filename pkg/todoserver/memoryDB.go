package todoserver

import "errors"

// MemoryDB is an in memory Database for the TodoServer
type MemoryDB struct {
	todos  []TodoEntry
	lastID int
}

func (db *MemoryDB) nextID() int {
	db.lastID++
	return db.lastID
}

// NewMemoryDB instantiates a new in memory database
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		[]TodoEntry{},
		1,
	}
}

// GetTodoByID ...
func (db *MemoryDB) GetTodoByID(id int) *TodoEntry {
	// Try to find the todo
	for _, v := range db.todos {
		if v.ID == id {
			return &v
		}
	}
	return nil
}

// GetTodosForUser ...
func (db *MemoryDB) GetTodosForUser(id int) []TodoEntry {
	return db.todos
}

// InsertTodo ...
func (db *MemoryDB) InsertTodo(todo TodoEntry) (*TodoEntry, error) {
	todo.ID = db.nextID()
	db.todos = append(db.todos, todo)
	return &todo, nil
}

// DeleteTodo ...
func (db *MemoryDB) DeleteTodo(id int) error {
	for i, v := range db.todos {
		if v.ID == id {
			// Remove from slice
			db.todos = append(db.todos[:i], db.todos[i+1:]...)
			return nil
		}
	}
	return errors.New("Not found")
}
