package domain

import "errors"

var (
	ErrInvalidTweetContent = errors.New("invalid tweet content")
	ErrInvalidFollowAction = errors.New("invalid follow action")
	ErrAlreadyFollowing    = errors.New("already following")
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
		return ErrAlreadyFollowing
	}
	u.Following[userID] = true
	return nil
}
