package application

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pedro00627/urblog/domain"
	"github.com/pedro00627/urblog/infrastructure"
	"github.com/pedro00627/urblog/infrastructure/db"
	"github.com/pedro00627/urblog/infrastructure/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewCreateTweetUseCase(t *testing.T) {
	tweetRepo := &mocks.MockTweetRepository{}
	userRepo := &mocks.MockUserRepository{}
	queue := &mocks.MockQueue{}

	useCase := NewCreateTweetUseCase(tweetRepo, userRepo, queue)

	assert.NotNil(t, useCase)
	assert.Equal(t, tweetRepo, useCase.tweetRepo)
	assert.Equal(t, userRepo, useCase.userRepo)
	assert.Equal(t, queue, useCase.queue)
}

func TestCreateTweetUseCase_Execute(t *testing.T) {
	type fields struct {
		tweetRepo db.TweetRepository
		userRepo  db.UserRepository
		queue     infrastructure.Queue
	}
	type args struct {
		userID  string
		content string
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Tweet
		wantErr assert.ErrorAssertionFunc
		mocks   func(f fields)
	}{
		{
			name: "success",
			fields: fields{
				tweetRepo: mocks.NewMockTweetRepository(ctrl),
				userRepo:  mocks.NewMockUserRepository(ctrl),
				queue:     mocks.NewMockQueue(ctrl),
			},
			args: args{
				userID:  "user1",
				content: "Hello, world!",
			},
			want: &domain.Tweet{
				UserID:    "user1",
				Content:   "Hello, world!",
				Timestamp: time.Now(),
			},
			wantErr: assert.NoError,
			mocks: func(f fields) {
				f.userRepo.(*mocks.MockUserRepository).EXPECT().FindByID(gomock.Eq("user1")).Return(domain.NewUser("user1", "User 1"), nil).Times(1)
				f.tweetRepo.(*mocks.MockTweetRepository).EXPECT().Save(gomock.Any()).Times(1)
				f.queue.(*mocks.MockQueue).EXPECT().WriteMessage(gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name: "user not found",
			fields: fields{
				tweetRepo: mocks.NewMockTweetRepository(ctrl),
				userRepo:  mocks.NewMockUserRepository(ctrl),
				queue:     mocks.NewMockQueue(ctrl),
			},
			args: args{
				userID:  "user1",
				content: "Hello, world!",
			},
			want:    nil,
			wantErr: assert.Error,
			mocks: func(f fields) {
				f.userRepo.(*mocks.MockUserRepository).EXPECT().FindByID(gomock.Eq("user1")).Return(nil, domain.ErrUserNotFound).Times(1)
			},
		},
		{
			name: "user not found not errpr",
			fields: fields{
				tweetRepo: mocks.NewMockTweetRepository(ctrl),
				userRepo:  mocks.NewMockUserRepository(ctrl),
				queue:     mocks.NewMockQueue(ctrl),
			},
			args: args{
				userID:  "user1",
				content: "Hello, world!",
			},
			want:    nil,
			wantErr: assert.Error,
			mocks: func(f fields) {
				f.userRepo.(*mocks.MockUserRepository).EXPECT().FindByID(gomock.Eq("user1")).Return(nil, nil).Times(1)
			},
		},
		{
			name: "error saving tweet",
			fields: fields{
				tweetRepo: mocks.NewMockTweetRepository(ctrl),
				userRepo:  mocks.NewMockUserRepository(ctrl),
				queue:     mocks.NewMockQueue(ctrl),
			},
			args: args{
				userID:  "user1",
				content: "Hello, world!",
			},
			want:    nil,
			wantErr: assert.Error,
			mocks: func(f fields) {
				f.userRepo.(*mocks.MockUserRepository).EXPECT().FindByID(gomock.Eq("user1")).Return(domain.NewUser("user1", "User 1"), nil).Times(1)
				f.tweetRepo.(*mocks.MockTweetRepository).EXPECT().Save(gomock.Any()).Return(errors.New("error saving tweet")).Times(1)
			},
		},
		{
			name: "error writing to queue",
			fields: fields{
				tweetRepo: mocks.NewMockTweetRepository(ctrl),
				userRepo:  mocks.NewMockUserRepository(ctrl),
				queue:     mocks.NewMockQueue(ctrl),
			},
			args: args{
				userID:  "user1",
				content: "Hello, world!",
			},
			want:    nil,
			wantErr: assert.Error,
			mocks: func(f fields) {
				f.userRepo.(*mocks.MockUserRepository).EXPECT().FindByID(gomock.Eq("user1")).Return(domain.NewUser("user1", "User 1"), nil).Times(1)
				f.tweetRepo.(*mocks.MockTweetRepository).EXPECT().Save(gomock.Any()).Times(1)
				f.queue.(*mocks.MockQueue).EXPECT().WriteMessage(gomock.Any()).Return(errors.New("error writing to queue")).Times(1)
			},
		},
		{
			name: "invalid tweet content",
			fields: fields{
				tweetRepo: mocks.NewMockTweetRepository(ctrl),
				userRepo:  mocks.NewMockUserRepository(ctrl),
				queue:     mocks.NewMockQueue(ctrl),
			},
			args: args{
				userID:  "user1",
				content: "",
			},
			want:    nil,
			wantErr: assert.Error,
			mocks: func(f fields) {
				f.userRepo.(*mocks.MockUserRepository).EXPECT().FindByID(gomock.Eq("user1")).Return(domain.NewUser("user1", "User 1"), nil).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &CreateTweetUseCase{
				tweetRepo: tt.fields.tweetRepo,
				userRepo:  tt.fields.userRepo,
				queue:     tt.fields.queue,
			}
			tt.mocks(tt.fields)
			got, err := uc.Execute(tt.args.userID, tt.args.content)
			if !tt.wantErr(t, err, fmt.Sprintf("Execute(%v, %v)", tt.args.userID, tt.args.content)) {
				return
			}
			if tt.want != nil {
				assert.NotEmpty(t, got.ID, "Tweet ID should not be empty")
				assert.Equalf(t, tt.want.UserID, got.UserID, "Execute(%v, %v)", tt.args.userID, tt.args.content)
				assert.Equalf(t, tt.want.Content, got.Content, "Execute(%v, %v)", tt.args.userID, tt.args.content)
				assert.WithinDuration(t, tt.want.Timestamp, got.Timestamp, time.Second, "Execute(%v, %v)", tt.args.userID, tt.args.content)
			}
		})
	}
}
