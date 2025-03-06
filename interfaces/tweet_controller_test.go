package interfaces

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/pedro00627/urblog/application"
	"github.com/pedro00627/urblog/application/mocks"
	"github.com/pedro00627/urblog/domain"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewTweetController(t *testing.T) {
	type args struct {
		createTweet application.CreateTweet
	}
	tests := []struct {
		name string
		args args
		want *TweetController
	}{
		{
			name: "valid createTweet",
			args: args{
				createTweet: &application.CreateTweetUseCase{},
			},
			want: &TweetController{
				createTweet: &application.CreateTweetUseCase{},
			},
		},
		{
			name: "nil createTweet",
			args: args{
				createTweet: nil,
			},
			want: &TweetController{
				createTweet: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewTweetController(tt.args.createTweet), "NewTweetController(%v)", tt.args.createTweet)
		})
	}
}

func TestTweetController_CreateTweet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		createTweet *mocks.MockCreateTweet
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		setupMocks func(f *fields)
		wantStatus int
		wantBody   string
	}{
		{
			name: "valid request",
			fields: fields{
				createTweet: mocks.NewMockCreateTweet(ctrl),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/tweets", strings.NewReader(`{"user_id":"user1","content":"Hello, world!"}`)),
			},
			setupMocks: func(f *fields) {
				f.createTweet.EXPECT().Execute("user1", "Hello, world!").Return(&domain.Tweet{
					ID:        "tweet1",
					UserID:    "user1",
					Content:   "Hello, world!",
					Timestamp: time.Date(2023, 10, 10, 10, 0, 0, 0, time.UTC),
				}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"id":"tweet1","user_id":"user1","content":"Hello, world!","timestamp":"2023-10-10 10:00:00 +0000 UTC"}`,
		},
		{
			name: "invalid request body",
			fields: fields{
				createTweet: mocks.NewMockCreateTweet(ctrl),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/tweets", strings.NewReader(`invalid json`)),
			},
			setupMocks: func(f *fields) {},
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid request body\n",
		},
		{
			name: "error creating tweet",
			fields: fields{
				createTweet: mocks.NewMockCreateTweet(ctrl),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/tweets", strings.NewReader(`{"user_id":"user1","content":"Hello, world!"}`)),
			},
			setupMocks: func(f *fields) {
				f.createTweet.EXPECT().Execute("user1", "Hello, world!").Return(nil, errors.New("error creating tweet"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "error creating tweet\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(&tt.fields)
			c := &TweetController{
				createTweet: tt.fields.createTweet,
			}
			c.CreateTweet(tt.args.w, tt.args.r)
			res := tt.args.w.(*httptest.ResponseRecorder)
			assert.Equal(t, tt.wantStatus, res.Code)
			if tt.wantStatus == http.StatusOK {
				assert.JSONEq(t, tt.wantBody, res.Body.String())
			} else {
				assert.Equal(t, tt.wantBody, res.Body.String())
			}
		})
	}
}
