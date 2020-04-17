package mongoutil

import (
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ParseID takes an interface which represents an ObjectID
func ParseID(id interface{}) *primitive.ObjectID {
	var objID primitive.ObjectID
	// Parse to objectID
	switch i := id.(type) {
	case string:
		objID, _ = primitive.ObjectIDFromHex(i)
		break
	case primitive.ObjectID:
		objID = i
		break
	default:
		log.Println("Received invalid id type")
		return nil
	}

	return &objID
}
