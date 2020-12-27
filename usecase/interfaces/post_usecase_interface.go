package interfaces

import "github.com/timofef/technopark_subd_sem_project/models"

type PostUsecase interface {
	GetPostDetails(id *string, related []byte) (*models.PostFull, error)
	EditPost(id *string, update *models.PostUpdate) (*models.Post, error)
}