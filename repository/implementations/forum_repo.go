package implementations

import (
	"github.com/jackc/pgx"
	"github.com/timofef/technopark_subd_sem_project/repository/interfaces"
)

type ForumRepo struct {
	db *pgx.ConnPool
}

func NewForumRepo(pool *pgx.ConnPool) interfaces.ForumRepository {
	return &ForumRepo{db: pool}
}

func (f *ForumRepo) CreateForum() {
	panic("implement me")
}

func (f *ForumRepo) GetDetailsBySlug() {
	panic("implement me")
}

func (f *ForumRepo) CreateBranchBySlug() {
	panic("implement me")
}

func (f *ForumRepo) GetUsersBySlug() {
	panic("implement me")
}

func (f *ForumRepo) GetTreadsBySlug() {
	panic("implement me")
}

func (f *ForumRepo) PrepareStatements() {
	panic("implement me")
}