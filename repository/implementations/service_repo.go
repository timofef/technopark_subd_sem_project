package implementations

import (
	"github.com/jackc/pgx"
	"github.com/timofef/technopark_subd_sem_project/models"
	"github.com/timofef/technopark_subd_sem_project/repository/interfaces"
)

type ServiceRepo struct {
	db *pgx.ConnPool
}

func NewServiceRepo(pool *pgx.ConnPool) interfaces.ServiceRepository {
	newRepo := &ServiceRepo{db: pool}

	return newRepo
}

func (s *ServiceRepo) GetStatus() (*models.Status, error) {
	tx, err := s.db.Begin()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	status := models.Status{}
	tx.QueryRow("SELECT " +
		"(SELECT count(id) FROM forums) as forums, "+
		"(SELECT count(id) FROM posts) as posts, "+
		"(SELECT count(id) FROM users) as users, "+
		"(SELECT count(id) FROM threads) as threads").
		Scan(&status.Forum,
			&status.Post,
			&status.User,
			&status.Thread)

	return &status, nil
}

func (s *ServiceRepo) Clear() error {
	tx, err := s.db.Begin()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	tx.Exec("TRUNCATE users, forums, threads, posts, votes, forum_users")

	return nil
}
