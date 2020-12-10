package implementations

import (
	usecase "github.com/timofef/technopark_subd_sem_project/usecase/interfaces"
	repository "github.com/timofef/technopark_subd_sem_project/repository/interfaces"
)

type UserUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userR repository.UserRepository) usecase.UserUsecase {
	return &UserUsecase{userRepo: userR}
}

func (uu * UserUsecase) CreateUser()  {

}

func (uu * UserUsecase) GetUser() {

}

func (uu * UserUsecase) UpdateUser() {

}