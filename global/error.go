package global

import "errors"

var ErrOrgExists = errors.New("organization already exists")
var ErrUserExists = errors.New("user already exists")

func CreateDuplicateErr(r string) (error) {
	return errors.New(r+" already exist")
}