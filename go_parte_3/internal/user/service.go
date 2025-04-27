package user

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

// Service provides high-level user management operations on a LocalStorage backend.
type Service struct {
	// storage is the underlying persistence for User entities.
	storage Storage

	// logger is our observability component to log.
	logger *zap.Logger
}

// NewService creates a new Service.
func NewService(storage Storage, logger *zap.Logger) *Service {
	if logger == nil {
		logger, _ = zap.NewProduction()
		defer logger.Sync() // flushes buffer, if any
	}
	
	return &Service{
		storage: storage,
		logger:  logger,
	}
}

// Create adds a brand-new user to the system.
// It sets CreatedAt and UpdatedAt to the current time and initializes Version to 1.
// Returns ErrEmptyID if user.ID is empty.
func (s *Service) Create(user *User) error {
	user.ID = uuid.NewString()
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.Version = 1

	if err := s.storage.Set(user); err != nil {
		s.logger.Error("failed to set user", zap.Error(err), zap.Any("user", user))
		return err
	}

	return nil
}

// Get retrieves a user by its ID.
// Returns ErrNotFound if no user exists with the given ID.
func (s *Service) Get(id string) (*User, error) {
	return s.storage.Read(id)
}

// Update modifies an existing user's data.
// It updates Name, Address, NickName, sets UpdatedAt to now and increments Version.
// Returns ErrNotFound if the user does not exist, or ErrEmptyID if user.ID is empty.
func (s *Service) Update(id string, user *UpdateFields) (*User, error) {
	existing, err := s.storage.Read(id)
	if err != nil {
		return nil, err
	}

	if user.Name != nil {
		existing.Name = *user.Name
	}

	if user.Address != nil {
		existing.Address = *user.Address
	}

	if user.NickName != nil {
		existing.NickName = *user.NickName
	}

	existing.UpdatedAt = time.Now()
	existing.Version++

	if err := s.storage.Set(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// Delete removes a user from the system by its ID.
// Returns ErrNotFound if the user does not exist.
func (s *Service) Delete(id string) error {
	return s.storage.Delete(id)
}
