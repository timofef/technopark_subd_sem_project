package implementations

import (
	"github.com/timofef/technopark_subd_sem_project/models"
	repository "github.com/timofef/technopark_subd_sem_project/repository/interfaces"
	usecase "github.com/timofef/technopark_subd_sem_project/usecase/interfaces"
)

type UserUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userR repository.UserRepository) usecase.UserUsecase {
	return &UserUsecase{userRepo: userR}
}

func (uu * UserUsecase) CreateUser(user *models.User, nickname string) (*models.Users, error)  {
	users, err := uu.userRepo.CreateUser(user, nickname)

	return users, err
}

func (uu * UserUsecase) GetUser(nickname string) (*models.User, error) {
	user, err := uu.userRepo.GetUserByNickname(nickname)

	return user, err
}

func (uu * UserUsecase) UpdateUser(newInfo *models.UserUpdate, nickname string) (*models.User, error) {
	if _, err := uu.userRepo.GetUserByNickname(nickname); err != nil {
		return &models.User{}, models.UserNotExists
	}

	user, err := uu.userRepo.UpdateUserByNickname(newInfo, nickname)

	return user, err
}