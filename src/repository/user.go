package repository

import (
	"crud-example-go/src/dto"
	"crud-example-go/src/errors"
	"crud-example-go/src/intreface/db"
	"github.com/jinzhu/copier"
)

type UserRepository struct {
	DB db.DB
}

type Tag struct {
	ID     *int `gorm:"primaryKey"`
	Name   string
	UserID *int
}

type User struct {
	ID   *int `gorm:"primaryKey"`
	Name string
	Tags []Tag `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (User) TableName() string {
	return "user"
}

func (Tag) TableName() string {
	return "tag"
}

func (ur *UserRepository) GetUsers() (*[]dto.User, error) {
	var users []User

	if result := ur.DB.Preload("Tags").Find(&users); result.GetError() != nil {
		return nil, result.GetError()
	} else {
		var dtoUsers []dto.User
		if err := copier.Copy(&dtoUsers, &users); err != nil {
			return nil, err
		}
		return &dtoUsers, nil
	}
}

func (ur *UserRepository) GetUser(id int) (*dto.User, error) {
	var user User
	if result := ur.DB.Preload("Tags").First(&user, id); result.GetError() != nil {
		if result.GetError() == errors.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.GetError()
	} else {
		return userModelToDto(&user)
	}
}

func (ur *UserRepository) SaveUser(userDto *dto.User) (*dto.User, error) {
	var user User
	if err := copier.Copy(&user, userDto); err != nil {
		return nil, err
	}

	if err := ur.clearTagsByUserId(user.ID); err != nil {
		return nil, err
	}

	if result := ur.DB.Save(&user); result.GetError() != nil {
		return nil, result.GetError()
	} else {
		return userModelToDto(&user)
	}
}

func (ur *UserRepository) DeleteUser(id int) error {
	if result := ur.DB.Delete(&User{}, id); result.GetError() != nil {
		return result.GetError()
	}
	return nil
}

func (ur *UserRepository) clearTagsByUserId(id *int) error {
	if id == nil {
		return nil
	}

	var tags []Tag

	if result := ur.DB.Find(&tags, "user_id = ?", id); result.GetError() != nil {
		return result.GetError()
	}
	for _, tag := range tags {
		if result := ur.DB.Delete(&tag); result.GetError() != nil {
			return result.GetError()
		}
	}
	return nil
}

func userModelToDto(user *User) (*dto.User, error) {
	var dtoUser dto.User
	if err := copier.Copy(&dtoUser, user); err != nil {
		return nil, err
	}
	return &dtoUser, nil
}
