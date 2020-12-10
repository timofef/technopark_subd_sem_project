package interfaces

import "github.com/timofef/technopark_subd_sem_project/models"

type UserRepository interface {
	CreateUser(user *models.User, nickname string) (models.Users, error)
	GetUserByNickname(nickname string) (models.User, error)
	UpdateUserByNickname(newInfo *models.User, nickname string) (models.User, error)
	PrepareStatements() error
}
