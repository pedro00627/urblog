package application

import (
	"testing"

	"github.com/pedro00627/urblog/domain"
	"github.com/stretchr/testify/assert"
)

type MockUserRepository struct {
	users map[string]*domain.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *MockUserRepository) Save(user *domain.User) error {
	r.users[user.ID] = user
	return nil
}

func (r *MockUserRepository) FindByID(userID string) (*domain.User, error) {
	user, exists := r.users[userID]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (r *MockUserRepository) FindByName(s string) (*domain.User, error) {
	panic("unimplemented")
}

func TestFollowUser(t *testing.T) {
	userRepo := NewMockUserRepository()
	kafkaWriter := &MockKafkaWriter{}
	followUserUseCase := NewFollowUserUseCase(userRepo, kafkaWriter)

	user := domain.NewUser("user1", "User One")
	userRepo.Save(user)

	err := followUserUseCase.Execute("user1", "user2")
	assert.NoError(t, err)
	assert.True(t, user.Following["user2"])
}
