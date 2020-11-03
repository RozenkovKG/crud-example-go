package repository

import "crud-example-go/src/dto"

type UserRepository interface {
	GetUsers() (*[]dto.User, error)
	GetUser(id int) (*dto.User, error)
	SaveUser(user *dto.User) (*dto.User, error)
	DeleteUser(id int) error
}
