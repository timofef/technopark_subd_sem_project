package interfaces

import "github.com/timofef/technopark_subd_sem_project/models"

type ThreadRepository interface {
	CreateThread(thread *models.Thread) (*models.Thread, error)
	GetThreadBySlug(slug string) (*models.Thread, error)
	GetThreadBySlugOrId(slugOrId interface{}) (*models.Thread, error)
	VoteForThread(thread *models.Thread, voice *models.Vote) (*models.Thread, error)
	GetThreadPosts(thread *models.Thread, limit, since, sort , desc []byte) (*models.Posts, error)
	GetThreadPostsFlat(thread *models.Thread, limit, since, desc []byte) (*models.Posts, error)
	GetThreadPostsTree(thread *models.Thread, limit, since, desc []byte) (*models.Posts, error)
	GetThreadPostsParentTree(thread *models.Thread, limit, since, desc []byte) (*models.Posts, error)
	UpdateThreadById(id int32, threadUpdate *models.ThreadUpdate) (error)
	PrepareStatements() error
}
