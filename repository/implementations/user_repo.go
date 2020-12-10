package implementations

import (
	"github.com/jackc/pgx"
	"github.com/timofef/technopark_subd_sem_project/repository/interfaces"
)

type UserRepo struct {
	db *pgx.ConnPool
}

func NewUserRepo(pool *pgx.ConnPool) interfaces.ForumRepository {
	return &ForumRepo{db: pool}
}

func (u *UserRepo) CreateUser() {
	panic("implement me")
}

func (u *UserRepo) GetUserByNickname() {
	panic("implement me")
}

func (u *UserRepo) UpdateUserByNickname() {
	panic("implement me")
}

func (u *UserRepo) PrepareStatements() {
	panic("implement me")
}
