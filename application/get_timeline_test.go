package application

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pedro00627/urblog/domain"
	"github.com/pedro00627/urblog/infrastructure/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetTimelineUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTweetRepo := mocks.NewMockTweetRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	useCase := NewGetTimelineUseCase(mockTweetRepo, mockUserRepo)

	inputDate := time.Date(2025, 3, 4, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		userID     string
		limit      int
		offset     int
		setupMocks func()
		wantTweets []*domain.Tweet
		wantErr    error
	}{
		{
			name:   "success",
			userID: "user1",
			limit:  10,
			offset: 0,
			setupMocks: func() {
				user := &domain.User{
					ID:       "user1",
					Username: "user1",
					Following: map[string]bool{
						"user2": true,
					},
				}
				tweet1 := &domain.Tweet{
					ID:        "tweet1",
					UserID:    "user2",
					Content:   "Hello, world!",
					Timestamp: inputDate.Add(-1 * time.Hour),
				}
				tweet2 := &domain.Tweet{
					ID:        "tweet2",
					UserID:    "user2",
					Content:   "Another tweet",
					Timestamp: inputDate.Add(-2 * time.Hour),
				}

				mockUserRepo.EXPECT().FindByID("user1").Return(user, nil)
				mockUserRepo.EXPECT().FindByName("user2").Return(&domain.User{ID: "user2"}, nil)
				mockTweetRepo.EXPECT().FindByUserID("user2", 10, 0).Return([]*domain.Tweet{tweet1, tweet2}, nil)
			},
			wantTweets: []*domain.Tweet{
				{
					ID:        "tweet1",
					UserID:    "user2",
					Content:   "Hello, world!",
					Timestamp: inputDate.Add(-1 * time.Hour),
				},
				{
					ID:        "tweet2",
					UserID:    "user2",
					Content:   "Another tweet",
					Timestamp: inputDate.Add(-2 * time.Hour),
				},
			},
			wantErr: nil,
		},
		{
			name:   "user not found",
			userID: "user1",
			limit:  10,
			offset: 0,
			setupMocks: func() {
				mockUserRepo.EXPECT().FindByID("user1").Return(nil, domain.ErrUserNotFound)
			},
			wantTweets: nil,
			wantErr:    domain.ErrUserNotFound,
		},
		{
			name:   "error finding followed user",
			userID: "user1",
			limit:  10,
			offset: 0,
			setupMocks: func() {
				user := &domain.User{
					ID:       "user1",
					Username: "user1",
					Following: map[string]bool{
						"user2": true,
					},
				}
				mockUserRepo.EXPECT().FindByID("user1").Return(user, nil)
				mockUserRepo.EXPECT().FindByName("user2").Return(nil, errors.New("error finding user"))
			},
			wantTweets: nil,
			wantErr:    nil, // No error is returned in this case, just logs
		},
		{
			name:   "error finding tweets",
			userID: "user1",
			limit:  10,
			offset: 0,
			setupMocks: func() {
				user := &domain.User{
					ID:       "user1",
					Username: "user1",
					Following: map[string]bool{
						"user2": true,
					},
				}
				mockUserRepo.EXPECT().FindByID("user1").Return(user, nil)
				mockUserRepo.EXPECT().FindByName("user2").Return(&domain.User{ID: "user2"}, nil)
				mockTweetRepo.EXPECT().FindByUserID("user2", 10, 0).Return(nil, errors.New("error finding tweets"))
			},
			wantTweets: nil,
			wantErr:    errors.New("error finding tweets"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			tweets, err := useCase.Execute(tt.userID, tt.limit, tt.offset)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantTweets, tweets)
			}
		})
	}
}

func Test_paginateTweets(t *testing.T) {
	type args struct {
		offset    int
		limit     int
		allTweets []*domain.Tweet
	}
	tests := []struct {
		name string
		args args
		want []*domain.Tweet
	}{
		{
			name: "no tweets",
			args: args{
				offset:    0,
				limit:     10,
				allTweets: []*domain.Tweet{},
			},
			want: []*domain.Tweet{},
		},
		{
			name: "offset beyond length",
			args: args{
				offset: 10,
				limit:  10,
				allTweets: []*domain.Tweet{
					{ID: "tweet1", Content: "Hello"},
				},
			},
			want: []*domain.Tweet{},
		},
		{
			name: "limit beyond length",
			args: args{
				offset: 0,
				limit:  10,
				allTweets: []*domain.Tweet{
					{ID: "tweet1", Content: "Hello"},
					{ID: "tweet2", Content: "World"},
				},
			},
			want: []*domain.Tweet{
				{ID: "tweet1", Content: "Hello"},
				{ID: "tweet2", Content: "World"},
			},
		},
		{
			name: "within range",
			args: args{
				offset: 0,
				limit:  1,
				allTweets: []*domain.Tweet{
					{ID: "tweet1", Content: "Hello"},
					{ID: "tweet2", Content: "World"},
				},
			},
			want: []*domain.Tweet{
				{ID: "tweet1", Content: "Hello"},
			},
		},
		{
			name: "offset and limit within range",
			args: args{
				offset: 1,
				limit:  1,
				allTweets: []*domain.Tweet{
					{ID: "tweet1", Content: "Hello"},
					{ID: "tweet2", Content: "World"},
				},
			},
			want: []*domain.Tweet{
				{ID: "tweet2", Content: "World"},
			},
		},
		{
			name: "offset equals length",
			args: args{
				offset: 2,
				limit:  1,
				allTweets: []*domain.Tweet{
					{ID: "tweet1", Content: "Hello"},
					{ID: "tweet2", Content: "World"},
				},
			},
			want: []*domain.Tweet{},
		},
		{
			name: "end beyond length",
			args: args{
				offset: 1,
				limit:  2,
				allTweets: []*domain.Tweet{
					{ID: "tweet1", Content: "Hello"},
					{ID: "tweet2", Content: "World"},
				},
			},
			want: []*domain.Tweet{
				{ID: "tweet2", Content: "World"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, paginateTweets(tt.args.offset, tt.args.limit, tt.args.allTweets), "paginateTweets(%v, %v, %v)", tt.args.offset, tt.args.limit, tt.args.allTweets)
		})
	}
}

func Test_sortTweetsByNewest(t *testing.T) {
	inputDate := time.Date(2025, 3, 4, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		allTweets []*domain.Tweet
		want      []*domain.Tweet
	}{
		{
			name: "sort by newest",
			allTweets: []*domain.Tweet{
				{ID: "tweet1", Timestamp: inputDate.Add(-2 * time.Hour)},
				{ID: "tweet2", Timestamp: inputDate.Add(-1 * time.Hour)},
				{ID: "tweet3", Timestamp: inputDate},
			},
			want: []*domain.Tweet{
				{ID: "tweet3", Timestamp: inputDate},
				{ID: "tweet2", Timestamp: inputDate.Add(-1 * time.Hour)},
				{ID: "tweet1", Timestamp: inputDate.Add(-2 * time.Hour)},
			},
		},
		{
			name: "already sorted",
			allTweets: []*domain.Tweet{
				{ID: "tweet3", Timestamp: inputDate},
				{ID: "tweet2", Timestamp: inputDate.Add(-1 * time.Hour)},
				{ID: "tweet1", Timestamp: inputDate.Add(-2 * time.Hour)},
			},
			want: []*domain.Tweet{
				{ID: "tweet3", Timestamp: inputDate},
				{ID: "tweet2", Timestamp: inputDate.Add(-1 * time.Hour)},
				{ID: "tweet1", Timestamp: inputDate.Add(-2 * time.Hour)},
			},
		},
		{
			name: "reverse order",
			allTweets: []*domain.Tweet{
				{ID: "tweet1", Timestamp: inputDate.Add(-2 * time.Hour)},
				{ID: "tweet2", Timestamp: inputDate.Add(-1 * time.Hour)},
				{ID: "tweet3", Timestamp: inputDate},
			},
			want: []*domain.Tweet{
				{ID: "tweet3", Timestamp: inputDate},
				{ID: "tweet2", Timestamp: inputDate.Add(-1 * time.Hour)},
				{ID: "tweet1", Timestamp: inputDate.Add(-2 * time.Hour)},
			},
		},
		{
			name:      "empty list",
			allTweets: []*domain.Tweet{},
			want:      []*domain.Tweet{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortTweetsByNewest(tt.allTweets)
			assert.Equal(t, tt.want, tt.allTweets)
		})
	}
}
