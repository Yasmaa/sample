package auth

import (

	"golang.org/x/crypto/bcrypt"
)

func Hash(data string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword string, password string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}