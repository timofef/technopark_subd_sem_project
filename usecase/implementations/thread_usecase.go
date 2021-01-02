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

func (u *ThreadUsecase) CreatePosts(slugOrId *interface{}, posts *models.Posts) (*models.Posts, error) {
	/*thread, err := u.threadRepo.GetThreadBySlugOrId(slugOrId)
	if err != nil {
		return nil, err
	}*/

	newPosts, err := u.postRepo.CreatePosts(slugOrId, posts)

	return newPosts, err
}

func (u *ThreadUsecase) GetThread(slugOrId *interface{}) (*models.Thread, error) {
	thread, err := u.threadRepo.GetThreadBySlugOrId(slugOrId)
/*	if err != nil {
		return nil, err
	}*/

	return thread, err
}

func (u *ThreadUsecase) UpdateThread(slug_or_id *interface{}, threadUpdate *models.ThreadUpdate) (*models.Thread, error) {
	/*thread, err := u.threadRepo.GetThreadBySlugOrId(slug_or_id)
	if err != nil {
		return nil, err
	}*/

	thread, err := u.threadRepo.UpdateThreadById(slug_or_id, threadUpdate)
	/*if err != nil {
		return &models.Thread{}, err
	}*/

	return thread, err
}

func (u *ThreadUsecase) GetPosts(slugOrId *interface{}, limit, since, sort, desc []byte) (*models.Posts, error) {
	posts, err := u.threadRepo.GetThreadPosts(slugOrId, limit, since, sort, desc)

	return posts, err
}

func (u *ThreadUsecase) VoteForThread(slugOrId *interface{}, voice *models.Vote) (*models.Thread, error) {
	/*existingThread, err := u.threadRepo.GetThreadBySlugOrId(slugOrId)
	if err != nil {
		return nil, models.ThreadNotExists
	}*/

	/*_, err = u.userRepo.GetUserByNickname(&voice.Nickname)
	if err != nil {
		return nil, models.UserNotExists
	}*/

	thread, err := u.threadRepo.VoteForThread(slugOrId, voice)

	return thread, err
}
