package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrSelfFriendship = errors.New("cannot add yourself as a friend")

type Friendship struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FriendID  uuid.UUID `json:"friend_id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewFriendship(userID, friendID uuid.UUID) (*Friendship, error) {
	if userID == friendID {
		return nil, ErrSelfFriendship
	}
	return &Friendship{
		ID:        uuid.New(),
		UserID:    userID,
		FriendID:  friendID,
		CreatedAt: time.Now().UTC(),
	}, nil
}
