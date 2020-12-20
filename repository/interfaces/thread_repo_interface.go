package interfaces

import "github.com/timofef/technopark_subd_sem_project/models"

type ThreadRepository interface {
	CreateThread(thread *models.Thread) (*models.Thread, error)
	GetThreadBySlug(slug string) (*models.Thread, error)
	GetThreadBySlugOrId(slugOrId interface{}) (*models.Thread, error)
	VoteForThread(thread *models.Thread, voice *models.Vote) (*models.Thread, error)
	PrepareStatements() error
}
