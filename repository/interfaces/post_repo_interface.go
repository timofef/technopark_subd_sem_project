package interfaces

import "github.com/timofef/technopark_subd_sem_project/models"

type PostRepository interface {
	CreatePosts(slugOrId *interface{}, posts *models.Posts) (*models.Posts, error)
	GetPostById(id *string) (*models.Post, error)
	EditPost(id *string, update *models.PostUpdate) (*models.Post, error)
	PrepareStatements() error
}
