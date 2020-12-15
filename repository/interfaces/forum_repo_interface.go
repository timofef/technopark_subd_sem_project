package interfaces

import "github.com/timofef/technopark_subd_sem_project/models"

type ForumRepository interface {
	CreateForum(forum *models.Forum) (*models.Forum, error)
	GetDetailsBySlug(slug string) (*models.Forum, error)
	GetUsersBySlug()
	GetTreadsBySlug()
	PrepareStatements() error
}
