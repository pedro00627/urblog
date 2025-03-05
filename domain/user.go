package domain

import "errors"

var (
	ErrInvalidTweetContent = errors.New("invalid tweet content")
	ErrInvalidFollowAction = errors.New("invalid follow action")
	ErrUserNotFound        = errors.New("user not found")
)

type User struct {
	ID        string
	Username  string
	Following map[string]bool
}

func NewUser(id, username string) *User {
	return &User{
		ID:        id,
		Username:  username,
		Following: make(map[string]bool),
	}
}

func (u *User) Follow(userID string) error {
	if u.ID == userID {
		return ErrInvalidFollowAction
	}
	if _, exists := u.Following[userID]; exists {
		return ErrInvalidFollowAction
	}
	u.Following[userID] = true
	return nil
}
