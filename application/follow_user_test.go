package application

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pedro00627/urblog/domain"
	"github.com/pedro00627/urblog/infrastructure/repositories/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFollowUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepo := mocks.NewMockUserRepository(ctrl)

	userRepo.EXPECT().Save(gomock.Any()).Times(2)
	userRepo.EXPECT().FindByID("user1").Return(domain.NewUser("user1", "user1"), nil).Times(1)
	userRepo.EXPECT().FindByID("user2").Return(domain.NewUser("user2", "user2"), nil).Times(1)
	kafkaWriter := &MockKafkaWriter{}
	followUserUseCase := NewFollowUserUseCase(userRepo, kafkaWriter)

	user := domain.NewUser("user1", "user1")
_:
	userRepo.Save(user)

	err := followUserUseCase.Execute("user1", "user2")
	assert.NoError(t, err)
}
