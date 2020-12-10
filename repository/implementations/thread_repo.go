package implementations

import (
	"github.com/jackc/pgx"
	"github.com/timofef/technopark_subd_sem_project/repository/interfaces"
)

type ThreadRepo struct {
	db *pgx.ConnPool
}

func NewThreadRepo(pool *pgx.ConnPool) interfaces.ThreadRepository {
	return &ThreadRepo{db: pool}
}

func (t * ThreadRepo) InsertNewPost() {

}

func (t * ThreadRepo) PrepareStatements() {

}