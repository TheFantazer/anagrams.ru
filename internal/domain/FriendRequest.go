package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSelfFriendRequest          = errors.New("cannot send friend request to yourself")
	ErrInvalidFriendRequestStatus = errors.New("invalid friend request status")
)

const (
	FriendRequestStatusPending  = "pending"
	FriendRequestStatusAccepted = "accepted"
	FriendRequestStatusRejected = "rejected"
)

type FriendRequest struct {
	ID         uuid.UUID `json:"id"`
	FromUserID uuid.UUID `json:"from_user_id"`
	ToUserID   uuid.UUID `json:"to_user_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewFriendRequest(fromUserId, toUserID uuid.UUID) (*FriendRequest, error) {
	if fromUserId == toUserID {
		return nil, ErrSelfFriendRequest
	}
	return &FriendRequest{
		ID:         uuid.New(),
		FromUserID: fromUserId,
		ToUserID:   toUserID,
		Status:     FriendRequestStatusPending,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}, nil
}

func (r *FriendRequest) Accept() {
	r.Status = FriendRequestStatusAccepted
	r.UpdatedAt = time.Now().UTC()
}

func (r *FriendRequest) Reject() {
	r.Status = FriendRequestStatusRejected
	r.UpdatedAt = time.Now().UTC()
}

func (r *FriendRequest) IsStatusPending() bool {
	return r.Status == FriendRequestStatusPending
}
