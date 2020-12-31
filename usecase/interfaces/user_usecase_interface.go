package interfaces

import "github.com/timofef/technopark_subd_sem_project/models"

type UserUsecase interface {
	CreateUser(user *models.User, nickname *string) (*models.Users, error)
	GetUser(nickname *string) (*models.User, error)
	UpdateUser(newInfo *models.UserUpdate, nickname *string) (*models.User, error)
}