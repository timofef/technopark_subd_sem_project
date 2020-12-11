package implementations

import (
	"github.com/timofef/technopark_subd_sem_project/models"
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

func (fu *ForumUsecase) CreateForum(forum *models.Forum) (models.Forum, error)  {

	if _, err := fu.userRepo.GetUserByNickname(forum.User); err != nil {
		return models.Forum{}, models.UserNotExists
	}
	newForum, err := fu.forumRepo.CreateForum(forum)

	return newForum, err
}

func (fu *ForumUsecase) CreateThread() {

}

func (fu *ForumUsecase) GetForumDetails() {

}

func (fu *ForumUsecase) GetForumThreads() {

}

func (fu *ForumUsecase) GetForumUsers() {

}
