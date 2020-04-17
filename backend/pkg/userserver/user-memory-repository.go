package userserver

import "errors"

// UserMemoryRepository is an in memory repository
type UserMemoryRepository struct {
	store  []User
	lastID int64
}

func (repo *UserMemoryRepository) nextID() int64 {
	repo.lastID++
	return repo.lastID
}

// NewUserMemoryRepository instantiates a new in memory database
func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		[]User{},
		1,
	}
}

// GetUserByEmail ...
func (repo *UserMemoryRepository) GetUserByEmail(email string) *User {
	for _, v := range repo.store {
		if v.Email == email {
			return &v
		}
	}
	return nil
}

// GetUserByID ...
func (repo *UserMemoryRepository) GetUserByID(id interface{}) *User {
	for _, v := range repo.store {
		if v.ID == id {
			return &v
		}
	}
	return nil
}

// InsertUser ...
func (repo *UserMemoryRepository) InsertUser(user User) (*User, error) {
	// Check if email is already used
	duplicate := repo.GetUserByEmail(user.Email)
	if duplicate != nil {
		return nil, errors.New("User with that email is already registered")
	}
	user.ID = repo.nextID()
	repo.store = append(repo.store, user)
	return &user, nil
}
