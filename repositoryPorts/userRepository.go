package repository

import "github.com/needsomesleeptd/annotater-core/models"

type IUserRepository interface {
	GetUserByLogin(login string) (*models.User, error)
	GetUserByID(id uint64) (*models.User, error)
	UpdateUserByLogin(login string, user *models.User) error
	DeleteUserByLogin(login string) error
	CreateUser(user *models.User) error
	GetAllUsers() ([]models.User, error)
}
