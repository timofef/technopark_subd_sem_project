package interfaces

import "github.com/timofef/technopark_subd_sem_project/models"

type ThreadUsecase interface {
	CreatePosts(slugOrId interface{}, posts *models.Posts) (*models.Posts, error)
	GetThread(slug_or_id interface{}) (*models.Thread, error)
	UpdateThread()
	GetPosts()
	VoteForThread(slugOrId interface{}, voice *models.Vote) (*models.Thread, error)
}