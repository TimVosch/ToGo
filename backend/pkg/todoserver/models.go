package todoserver

// TodoEntry is a todo by a user
type TodoEntry struct {
	ID      interface{} `json:"id" bson:"_id,omitempty"`
	Title   string      `json:"title"`
	OwnerID interface{} `json:"ownerID"`
}
