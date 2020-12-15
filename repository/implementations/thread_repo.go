package implementations

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/timofef/technopark_subd_sem_project/models"
	"github.com/timofef/technopark_subd_sem_project/repository/interfaces"
)

type ThreadRepo struct {
	db *pgx.ConnPool
}

func NewThreadRepo(pool *pgx.ConnPool) interfaces.ThreadRepository {
	new := &ThreadRepo{db: pool}

	if err := new.PrepareStatements(); err != nil {
		fmt.Println(err)
	}

	return new
}

func (t *ThreadRepo) CreateThread(thread *models.Thread) (*models.Thread, error) {
	tx, err := t.db.Begin()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	row := tx.QueryRow("insert_thread",
		thread.Author,
		thread.Forum,
		thread.Message,
		thread.Slug,
		thread.Title,
		thread.Created)

	if err = row.Scan(&thread.ID); err != nil {
		return nil, err
	}

	return thread, nil
}

func (t *ThreadRepo) GetThreadBySlug(slug string) (*models.Thread, error) {
	tx, err := t.db.Begin()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	rows := tx.QueryRow("get_thread_by_slug", slug)

	var thread = models.Thread{}
	err = rows.Scan(&thread.ID,
		&thread.Author,
		&thread.Created,
		&thread.Forum,
		&thread.Message,
		&thread.Title,
		&thread.Votes)

	if err != nil {
		return nil, models.ThreadNotExists
	}

	return &thread, nil
}

func (t * ThreadRepo) PrepareStatements() error {
	_, err := t.db.Prepare("insert_thread",
		"INSERT INTO threads (author, forum, message, slug, title, created) "+
			"VALUES ($1, $2, $3, CASE WHEN $4 = '' THEN NULL ELSE $4 END, $5, $6) "+
			"RETURNING id ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_thread_by_slug",
		"SELECT threads.id, threads.author, threads.created, threads.forum, "+
			"threads.message, threads.title, threads.votes "+
			"FROM threads "+
			"WHERE threads.slug = $1 ",
	)
	if err != nil {
		return err
	}

	return nil
}