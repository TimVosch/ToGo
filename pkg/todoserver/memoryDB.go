package todoserver

// MemoryDB is an in memory Database for the TodoServer
type MemoryDB struct {
	todos []TodoEntry
}

// NewMemoryDB instantiates a new in memory database
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		[]TodoEntry{},
	}
}

// GetTodosForUser ...
func (db *MemoryDB) GetTodosForUser(id int) []TodoEntry {
	return db.todos
}
