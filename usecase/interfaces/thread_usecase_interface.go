package interfaces

import "github.com/timofef/technopark_subd_sem_project/models"

type ThreadUsecase interface {
	CreatePosts(slugOrId *interface{}, posts *models.Posts) (*models.Posts, error)
	GetThread(slug_or_id *interface{}) (*models.Thread, error)
	UpdateThread(slug_or_id *interface{}, threadUpdate *models.ThreadUpdate) (*models.Thread, error)
	GetPosts(slug_or_id *interface{}, limit, since, sort, desc []byte) (*models.Posts, error)
	VoteForThread(slugOrId *interface{}, voice *models.Vote) (*models.Thread, error)
}