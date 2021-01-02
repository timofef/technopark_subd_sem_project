package interfaces

import (
	"github.com/jackc/pgx"
	"github.com/timofef/technopark_subd_sem_project/models"
)

type ThreadRepository interface {
	CreateThread(thread *models.Thread) (*models.Thread, error)
	GetThreadBySlug(slug string) (*models.Thread, error)
	GetThreadById(id int32) (*models.Thread, error)
	GetThreadBySlugOrId(slugOrId *interface{}) (*models.Thread, error)
	VoteForThread(slugOrId *interface{}, voice *models.Vote) (*models.Thread, error)
	GetThreadPosts(slugOrId *interface{}, limit, since, sort , desc []byte) (*models.Posts, error)
	GetThreadPostsFlat(tx *pgx.Tx, thread int32, limit, since, desc []byte) (*models.Posts, error)
	GetThreadPostsTree(tx *pgx.Tx, thread int32, limit, since, desc []byte) (*models.Posts, error)
	GetThreadPostsParentTree(tx *pgx.Tx, thread int32, limit, since, desc []byte) (*models.Posts, error)
	UpdateThreadById(slugOrId *interface{}, threadUpdate *models.ThreadUpdate) (*models.Thread, error)
	PrepareStatements() error
}
