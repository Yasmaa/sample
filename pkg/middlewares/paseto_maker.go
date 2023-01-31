package middlewares

import (
	"errors"
	"time"

	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       paseto.V2
	symmetricKey []byte
}

func (p *PasetoMaker) CreateToken(id uint,email string, duration time.Duration) (string, error) {
	payload, err := NewPayload(id ,email, duration)
	if err != nil {
		return "", nil
	}
	return p.paseto.Encrypt(p.symmetricKey, payload, nil)
}
func (p *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := p.paseto.Decrypt(token, p.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrorInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, ErrorInvalidToken
	}
	return payload, nil
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) < minLengthSecretKey {
		return nil, errors.New("invalid symmetric key length")
	}
	maker := &PasetoMaker{
		paseto:       *paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}