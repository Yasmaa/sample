package middlewares

import (
	"errors"
	"time"

	// "github.com/google/uuid"
)

var (
	ErrorInvalidToken = errors.New("invalid token")
	ErrorExpiredToken = errors.New("expired token")
)

// Payload contains data for the token
type Payload struct {
	ID        uint `json:"id"`
	Email  string    `json:"email"`
	IssueAt   time.Time `json:"issue_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(id uint,email string, duration time.Duration) (*Payload, error) {
	// tokenId, err := uuid.NewRandom()
	// if err != nil {
	// 	return nil, err
	// }
	now := time.Now()
	payload := &Payload{
		ID:        id,
		Email:  email,
		IssueAt:   now,
		ExpiredAt: now.Add(duration),
	}
	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrorExpiredToken
	}
	return nil
}