package implementations

import (
	"github.com/timofef/technopark_subd_sem_project/models"
	repository "github.com/timofef/technopark_subd_sem_project/repository/interfaces"
	usecase "github.com/timofef/technopark_subd_sem_project/usecase/interfaces"
)

type ForumUsecase struct {
	forumRepo  repository.ForumRepository
	userRepo   repository.UserRepository
	threadRepo repository.ThreadRepository
}

func NewForumUsecase(forumR repository.ForumRepository, userR repository.UserRepository, threadR repository.ThreadRepository) usecase.ForumUsecase {
	return &ForumUsecase{forumRepo: forumR, userRepo: userR, threadRepo: threadR}
}

func (fu *ForumUsecase) CreateForum(forum *models.Forum) (*models.Forum, error) {
	newForum, err := fu.forumRepo.CreateForum(forum)

	return newForum, err
}

func (fu *ForumUsecase) CreateThread(thread *models.Thread) (*models.Thread, error) {
	newThread, err := fu.threadRepo.CreateThread(thread)
	if err == nil {
		return newThread, nil
	} else if err != models.ForumNotExists {
		existingThread, err := fu.threadRepo.GetThreadBySlug(thread.Slug)
		if err == nil {
			return existingThread, models.ThreadExists
		}
	}

	return newThread, err
}

func (fu *ForumUsecase) GetForumDetails(slug *string) (*models.Forum, error) {
	forum, err := fu.forumRepo.GetDetailsBySlug(slug)

	return forum, err
}

func (fu *ForumUsecase) GetForumThreads(slug *string, since, desc, limit []byte) (*models.Threads, error) {
	threads, err := fu.forumRepo.GetThreads(slug, since, desc, limit)
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (fu *ForumUsecase) GetForumUsers(slug *string, since, desc, limit []byte) (*models.Users, error) {
	users, err := fu.forumRepo.GetUsersBySlug(slug, since, desc, limit)
	if err != nil {
		return nil, err
	}

	return users, nil
}
