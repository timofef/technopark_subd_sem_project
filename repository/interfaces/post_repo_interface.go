package interfaces

import "github.com/timofef/technopark_subd_sem_project/models"

type PostRepository interface {
	CreatePosts(posts *models.Posts, thread *models.Thread) (*models.Posts, error)
	PrepareStatements() error
}
