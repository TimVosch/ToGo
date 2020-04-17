package userserver

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserMongoRepository is an in memory repository
type UserMongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewUserMongoRepository instantiates a new in memory database
func NewUserMongoRepository(connectionString, database, collection string) *UserMongoRepository {
	r := &UserMongoRepository{}

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

// GetUserByEmail ...
func (r *UserMongoRepository) GetUserByEmail(email string) *User {
	var user User

	// Find by id and decode
	ctx, close := context.WithTimeout(context.Background(), time.Second*5)
	defer close()
	result := r.collection.FindOne(ctx, bson.M{"email": email})
	if err := result.Err(); err != nil {
		log.Println("Error occured while fetching user by email: ", err)
		return nil
	}

	err := result.Decode(&user)
	if err != nil {
		log.Println("Error occured while decoding user by email: ", err)
		return nil
	}

	return &user
}

// GetUserByID ...
func (r *UserMongoRepository) GetUserByID(id interface{}) *User {
	var user User

	// Find by id and decode
	ctx, close := context.WithTimeout(context.Background(), time.Second*5)
	defer close()
	result := r.collection.FindOne(ctx, bson.M{"_id": id})
	if err := result.Err(); err != nil {
		log.Println("Error occured while fetching user by ID: ", err)
		return nil
	}

	err := result.Decode(&user)
	if err != nil {
		log.Println("Error occured while decoding user by ID: ", err)
		return nil
	}

	return &user
}

// InsertUser ...
func (r *UserMongoRepository) InsertUser(user User) (*User, error) {
	ctx, close := context.WithTimeout(context.Background(), time.Second*5)
	defer close()
	res, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		log.Println("Error occured while inserting todo: ", err)
		return nil, err
	}

	return r.GetUserByID(res.InsertedID), nil
}
