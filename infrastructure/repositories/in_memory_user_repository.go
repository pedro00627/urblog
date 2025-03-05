package repositories

import "github.com/pedro00627/urblog/domain"

type InMemoryUserRepository struct {
	usersByID   map[string]*domain.User
	usersByName map[string]string
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		usersByID:   make(map[string]*domain.User),
		usersByName: make(map[string]string),
	}
}

func (r *InMemoryUserRepository) Save(user *domain.User) error {
	r.usersByID[user.ID] = user
	r.usersByName[user.Username] = user.ID
	return nil
}

func (r *InMemoryUserRepository) FindByID(userID string) (*domain.User, error) {
	user, exists := r.usersByID[userID]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (r *InMemoryUserRepository) FindByName(username string) (*domain.User, error) {
	userID, exists := r.usersByName[username]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	return r.FindByID(userID)
}
