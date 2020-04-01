package todoserver

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

// GetTodosForUser ...
func (db *MemoryDB) GetTodosForUser(id int) []TodoEntry {
	return db.todos
}

// InsertTodo ...
func (db *MemoryDB) InsertTodo(todo TodoEntry) error {
	todo.ID = db.nextID()
	db.todos = append(db.todos, todo)
	return nil
}
