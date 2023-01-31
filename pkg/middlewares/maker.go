package middlewares

import "time"

type Maker interface {
	CreateToken(id uint, email string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}