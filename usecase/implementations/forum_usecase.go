package implementations

import (
	repository "github.com/timofef/technopark_subd_sem_project/repository/interfaces"
	usecase "github.com/timofef/technopark_subd_sem_project/usecase/interfaces"
)

type ForumUsecase struct {
	forumRepo repository.ForumRepository
	userRepo repository.UserRepository
}

func NewForumUsecase(forumR repository.ForumRepository, userR repository.UserRepository) usecase.ForumUsecase {
	return &ForumUsecase{forumRepo: forumR, userRepo: userR}
}

func (fu * ForumUsecase) CreateForum()  {

}

func (fu * ForumUsecase) CreateThread() {

}

func (fu * ForumUsecase) GetForumDetails() {

}

func (fu * ForumUsecase) GetForumThreads() {

}

func (fu * ForumUsecase) GetForumUsers() {

}
