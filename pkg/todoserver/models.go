package todoserver

// TodoEntry is a todo by a user
type TodoEntry struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	OwnerID int64  `json:"ownerID"`
}
