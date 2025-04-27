package user

import "errors"

// ErrNotFound is returned when a user with the given ID is not found.
var ErrNotFound = errors.New("user not found")

// ErrEmptyID is returned when trying to store a user with an empty ID.
var ErrEmptyID = errors.New("empty user ID")

// Storage is the main interface for our storage layer.
type Storage interface {
	Set(user *User) error
	Read(id string) (*User, error)
	Delete(id string) error
}

// LocalStorage provides an in-memory implementation for storing users.
type LocalStorage struct {
	m map[string]*User
}

// NewLocalStorage instantiates a new LocalStorage with an empty map.
func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		m: map[string]*User{},
	}
}

// Set stores or updates a user in the local storage.
// Returns ErrEmptyID if the user has an empty ID.
func (l *LocalStorage) Set(user *User) error {
	if user.ID == "" {
		return ErrEmptyID
	}

	l.m[user.ID] = user
	return nil
}

// Read retrieves a user from the local storage by ID.
// Returns ErrNotFound if the user is not found.
func (l *LocalStorage) Read(id string) (*User, error) {
	u, ok := l.m[id]
	if !ok {
		return nil, ErrNotFound
	}

	return u, nil
}

// Delete removes a user from the local storage by ID.
// Returns ErrNotFound if the user does not exist.
func (l *LocalStorage) Delete(id string) error {
	_, err := l.Read(id)
	if err != nil {
		return err
	}

	delete(l.m, id)
	return nil
}
