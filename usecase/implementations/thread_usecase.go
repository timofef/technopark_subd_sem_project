package implementations

import (
	"github.com/timofef/technopark_subd_sem_project/models"
	repository "github.com/timofef/technopark_subd_sem_project/repository/interfaces"
	usecase "github.com/timofef/technopark_subd_sem_project/usecase/interfaces"
)

type ThreadUsecase struct {
	threadRepo repository.ThreadRepository
	forumRepo  repository.ForumRepository
	postRepo   repository.PostRepository
	userRepo   repository.UserRepository
}

func NewThreadUsecase(threadR repository.ThreadRepository,
	forumR repository.ForumRepository,
	postR repository.PostRepository,
	userR repository.UserRepository) usecase.ThreadUsecase {
	return &ThreadUsecase{
		threadRepo: threadR,
		forumRepo:  forumR,
		postRepo:   postR,
		userRepo:   userR}
}

func (u *ThreadUsecase) CreatePosts(slugOrId interface{}, posts *models.Posts) (*models.Posts, error) {
	thread, err := u.threadRepo.GetThreadBySlugOrId(slugOrId)
	if err != nil {
		return nil, err
	}

	newPosts, err := u.postRepo.CreatePosts(posts, thread)

	return newPosts, nil
}

func (u *ThreadUsecase) GetThread(slug_or_id interface{}) (*models.Thread, error) {
	thread, err := u.threadRepo.GetThreadBySlugOrId(slug_or_id)
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (u *ThreadUsecase) UpdateThread() {

}

func (u *ThreadUsecase) GetPosts() {

}

func (u *ThreadUsecase) VoteForThread(slugOrId interface{}, voice *models.Vote) (*models.Thread, error) {
	existingThread, err := u.threadRepo.GetThreadBySlugOrId(slugOrId)
	if err != nil {
		return nil, models.ThreadNotExists
	}
	thread, _ := u.threadRepo.GetThreadBySlug(existingThread.Slug)

	thread, err = u.threadRepo.VoteForThread(thread, voice)

	return thread, nil
}
