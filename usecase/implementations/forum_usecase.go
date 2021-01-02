package implementations

import (
	"github.com/timofef/technopark_subd_sem_project/models"
	repository "github.com/timofef/technopark_subd_sem_project/repository/interfaces"
	usecase "github.com/timofef/technopark_subd_sem_project/usecase/interfaces"
)

type ForumUsecase struct {
	forumRepo repository.ForumRepository
	userRepo repository.UserRepository
	threadRepo repository.ThreadRepository
}

func NewForumUsecase(forumR repository.ForumRepository, userR repository.UserRepository, threadR repository.ThreadRepository) usecase.ForumUsecase {
	return &ForumUsecase{forumRepo: forumR, userRepo: userR, threadRepo: threadR}
}

func (fu *ForumUsecase) CreateForum(forum *models.Forum) (*models.Forum, error) {
	/*user, err := fu.userRepo.GetUserByNickname(&forum.User)
	if err != nil {
		return nil, models.UserNotExists
	}
	forum.User = user.Nickname*/

	newForum, err := fu.forumRepo.CreateForum(forum)

	return newForum, err
}

func (fu *ForumUsecase) CreateThread(thread *models.Thread) (*models.Thread, error) {
	/*author, err := fu.userRepo.GetUserByNickname(&thread.Author)
	if err != nil {
		return nil, models.UserNotExists
	}
	thread.Author = author.Nickname

	forum, err := fu.forumRepo.GetDetailsBySlug(&thread.Forum)
	if err != nil {
		return nil, models.ForumNotExists
	}
	thread.Forum = forum.Slug*/

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
	/*forum, err := fu.forumRepo.GetDetailsBySlug(slug)
	if err != nil {
		return nil, models.ForumNotExists
	}*/

	threads, err := fu.forumRepo.GetThreads(slug, since, desc, limit)
	/*if err != nil {
		return nil, err
	}*/

	return threads, err
}

func (fu *ForumUsecase) GetForumUsers(slug *string, since, desc, limit []byte) (*models.Users, error) {
	/*forum, err := fu.forumRepo.GetDetailsBySlug(slug)
	if err != nil {
		return nil, models.ForumNotExists
	}*/

	users, err := fu.forumRepo.GetUsersBySlug(slug, since, desc, limit)
	/*if err != nil {
		return nil, err
	}*/

	return users, err
}
