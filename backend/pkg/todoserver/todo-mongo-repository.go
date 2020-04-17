package todoserver

import (
	"context"
	"log"
	"time"

	"github.com/timvosch/togo/pkg/common/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TodoMongoRepository is a repository for the TodoServer which implements mongodb
type TodoMongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewMongoRepository creates a new repository and connect to the mongo database
func NewMongoRepository(connectionString, database, collection string) *TodoMongoRepository {
	r := &TodoMongoRepository{}

	// Create mongo client
	clientOptions := options.Client().ApplyURI(connectionString).SetDirect(true)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalln("Could not create mongo client: ", err)
	}

	// Use client to connect to mongo server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln("Could not connect to mongo server: ", err)
	}

	// Check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalln("Could not Ping mongo server: ", err)
	}

	// Set repository
	r.client = client
	r.collection = client.Database(database).Collection(collection)

	return r
}

// GetTodoByID ...
func (r *TodoMongoRepository) GetTodoByID(id interface{}) *TodoEntry {
	var todo TodoEntry
	objID := mongoutil.ParseID(id)

	// Find by id and decode
	ctx, close := context.WithTimeout(context.Background(), time.Second*5)
	defer close()
	result := r.collection.FindOne(ctx, bson.M{"_id": objID})
	if err := result.Err(); err != nil {
		log.Println("Error occured while fetching todo by ID: ", err)
		return nil
	}

	err := result.Decode(&todo)
	if err != nil {
		log.Println("Error occured while decoding todo by ID: ", err)
		return nil
	}

	return &todo
}

// GetTodosForUser ...
func (r *TodoMongoRepository) GetTodosForUser(userID interface{}) []TodoEntry {
	todos := make([]TodoEntry, 0)
	objID := mongoutil.ParseID(userID)

	// Find by id and decode
	ctx, close := context.WithTimeout(context.Background(), time.Second*5)
	cur, err := r.collection.Find(ctx, bson.M{"ownerid": objID})
	defer close()
	if err != nil {
		log.Println("Error occured while fetching todos for user: ", err)
		return nil
	}

	// Loop results
	for cur.Next(ctx) {
		var todo TodoEntry
		err := cur.Decode(&todo)
		if err != nil {
			log.Println("Error occured while parsing todo for user: ", err)
			return nil
		}
		todos = append(todos, todo)
	}
	// Did the loop end because of an error?
	if err := cur.Err(); err != nil {
		log.Println("Error occured while parsing todos for user: ", err)
		return nil
	}

	return todos
}

// InsertTodo ...
func (r *TodoMongoRepository) InsertTodo(todo TodoEntry) (*TodoEntry, error) {
	todo.OwnerID = mongoutil.ParseID(todo.OwnerID)
	ctx, close := context.WithTimeout(context.Background(), time.Second*5)
	defer close()
	res, err := r.collection.InsertOne(ctx, todo)
	if err != nil {
		log.Println("Error occured while inserting todo: ", err)
		return nil, err
	}

	return r.GetTodoByID(res.InsertedID), nil
}

// DeleteTodo ...
func (r *TodoMongoRepository) DeleteTodo(id interface{}) error {
	objID := mongoutil.ParseID(id)

	ctx, close := context.WithTimeout(context.Background(), time.Second*5)
	defer close()
	_, err := r.collection.DeleteOne(ctx, bson.M{"ID": objID})
	if err != nil {
		log.Println("Error occured while inserting todo: ", err)
		return err
	}
	return nil
}
