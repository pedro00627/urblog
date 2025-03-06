package application

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pedro00627/urblog/domain"
	"github.com/pedro00627/urblog/infrastructure/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFollowUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	queue := mocks.NewMockQueue(ctrl)

	followUserUseCase := NewFollowUserUseCase(userRepo, queue)

	tests := []struct {
		name     string
		follower string
		followee string
		setup    func()
		wantErr  error
	}{
		{
			name:     "success",
			follower: "user1",
			followee: "user2",
			setup: func() {
				userRepo.EXPECT().FindByID("user1").Return(domain.NewUser("user1", "user1"), nil).Times(1)
				userRepo.EXPECT().FindByID("user2").Return(domain.NewUser("user2", "user2"), nil).Times(1)
				userRepo.EXPECT().Save(gomock.Any()).Times(1)
				queue.EXPECT().WriteMessage(gomock.Any()).Return(nil).Times(1)
			},
			wantErr: nil,
		},
		{
			name:     "follower not found",
			follower: "user1",
			followee: "user2",
			setup: func() {
				userRepo.EXPECT().FindByID("user1").Return(nil, domain.ErrUserNotFound).Times(1)
			},
			wantErr: domain.ErrUserNotFound,
		},
		{
			name:     "follower and followee are the same",
			follower: "user1",
			followee: "user1",
			setup: func() {
				userRepo.EXPECT().FindByID("user1").Return(domain.NewUser("user1", "user1"), nil).Times(2)
			},
			wantErr: domain.ErrInvalidFollowAction,
		},
		{
			name:     "followee not found",
			follower: "user1",
			followee: "user2",
			setup: func() {
				userRepo.EXPECT().FindByID("user1").Return(domain.NewUser("user1", "user1"), nil).Times(1)
				userRepo.EXPECT().FindByID("user2").Return(nil, domain.ErrUserNotFound).Times(1)
			},
			wantErr: domain.ErrUserNotFound,
		},
		{
			name:     "error saving user",
			follower: "user1",
			followee: "user2",
			setup: func() {
				userRepo.EXPECT().FindByID("user1").Return(domain.NewUser("user1", "user1"), nil).Times(1)
				userRepo.EXPECT().FindByID("user2").Return(domain.NewUser("user2", "user2"), nil).Times(1)
				userRepo.EXPECT().Save(gomock.Any()).Return(errors.New("error saving user")).Times(1)
			},
			wantErr: errors.New("error saving user"),
		},
		{
			name:     "error writing to queue",
			follower: "user1",
			followee: "user2",
			setup: func() {
				userRepo.EXPECT().FindByID("user1").Return(domain.NewUser("user1", "user1"), nil).Times(1)
				userRepo.EXPECT().FindByID("user2").Return(domain.NewUser("user2", "user2"), nil).Times(1)
				userRepo.EXPECT().Save(gomock.Any()).Times(1)
				queue.EXPECT().WriteMessage(gomock.Any()).Return(errors.New("error writing to queue")).Times(1)
			},
			wantErr: errors.New("error writing to queue"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := followUserUseCase.Execute(tt.follower, tt.followee)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
