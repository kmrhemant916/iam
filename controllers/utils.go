package controllers

import (
	"golang.org/x/crypto/bcrypt"
)

func GeneratehashedPassword(password []byte) (string, error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}