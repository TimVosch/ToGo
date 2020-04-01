package todoserver

// TodoEntry is a todo by a user
type TodoEntry struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	OwnerID int    `json:"ownerID"`
}
