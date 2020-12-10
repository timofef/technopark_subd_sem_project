package implementations

import (
	repository "github.com/timofef/technopark_subd_sem_project/repository/interfaces"
	usecase "github.com/timofef/technopark_subd_sem_project/usecase/interfaces"
)

type ThreadUsecase struct {
	threadRepo repository.ThreadRepository
	forumRepo repository.ForumRepository
}

func NewThreadUsecase(threadR repository.ThreadRepository, forumR repository.ForumRepository) usecase.ThreadUsecase {
	return &ThreadUsecase{threadRepo: threadR, forumRepo: forumR}
}

func (u *ThreadUsecase) CreatePost() {

}

func (u *ThreadUsecase) GetThread() {

}

func (u *ThreadUsecase) UpdateThread() {

}

func (u *ThreadUsecase) GetPosts() {

}

func (u *ThreadUsecase) VoteForThread() {

}