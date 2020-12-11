package interfaces

import "github.com/timofef/technopark_subd_sem_project/models"

type ForumRepository interface {
	CreateForum(forum *models.Forum) (models.Forum, error)
	GetDetailsBySlug()
	CreateBranchBySlug()
	GetUsersBySlug()
	GetTreadsBySlug()
	PrepareStatements() error
}
