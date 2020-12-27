package interfaces

import "github.com/timofef/technopark_subd_sem_project/models"

type ForumRepository interface {
	CreateForum(forum *models.Forum) (*models.Forum, error)
	GetDetailsBySlug(slug string) (*models.Forum, error)
	GetUsersBySlug(slug string, since, desc, limit []byte) (*models.Users, error)
	GetThreads(slug string, since, desc, limit []byte) (*models.Threads, error)
	PrepareStatements() error
}
