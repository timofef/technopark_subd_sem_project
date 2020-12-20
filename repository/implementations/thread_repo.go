package implementations

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/timofef/technopark_subd_sem_project/models"
	"github.com/timofef/technopark_subd_sem_project/repository/interfaces"
	"strconv"
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

	thread.Votes = 0

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
		&thread.Slug,
		&thread.Title,
		&thread.Votes,
	)

	if err != nil {
		return nil, models.ThreadNotExists
	}

	return &thread, nil
}

func (t *ThreadRepo) GetThreadBySlugOrId(slugOrId interface{}) (*models.Thread, error) {
	tx, err := t.db.Begin()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	var thread = models.Thread{}
	id, err := strconv.Atoi(slugOrId.(string))
	if err != nil {
		if err = tx.QueryRow("get_thread_by_slug", slugOrId).
			Scan(&thread.ID,
				&thread.Author,
				&thread.Created,
				&thread.Forum,
				&thread.Message,
				&thread.Slug,
				&thread.Title,
				&thread.Votes,
			); err != nil {
			return nil, models.ThreadNotExists
		}
	} else {
		if err = tx.QueryRow("get_thread_by_id", id).
			Scan(&thread.ID,
				&thread.Author,
				&thread.Created,
				&thread.Forum,
				&thread.Message,
				&thread.Slug,
				&thread.Title,
				&thread.Votes,
			); err != nil {
			return nil, models.ThreadNotExists
		}
	}

	return &thread, nil
}

func (t *ThreadRepo) VoteForThread(thread *models.Thread, voice *models.Vote) (*models.Thread, error) {
	tx, err := t.db.Begin()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	var alreadyVoted int32 = 0

	err = tx.QueryRow("check_vote",
		voice.Nickname,
		thread.ID,
	).Scan(&alreadyVoted)

	if err == nil && alreadyVoted == voice.Voice {
		return thread, models.SameVote
	} else {
		_, err = tx.Exec("insert_vote",
			voice.Voice,
			voice.Nickname,
			thread.ID,
		)
		if err != nil {
			_, err = tx.Exec("update_vote",
				voice.Voice,
				voice.Nickname,
				thread.ID,
			)
			thread.Votes += 2 * voice.Voice
		} else {
			thread.Votes += voice.Voice
		}
	}

	return thread, nil
}

func (t *ThreadRepo) PrepareStatements() error {
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
			"threads.message, threads.slug, threads.title, threads.votes "+
			"FROM threads "+
			"WHERE threads.slug = $1 ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_thread_by_id",
		"SELECT threads.id, threads.author, threads.created, threads.forum, "+
			"threads.message, threads.slug, threads.title, threads.votes "+
			"FROM threads "+
			"WHERE threads.id = $1 ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("check_thread_by_slug",
		"SELECT id, forum, slug FROM threads WHERE slug = $1")
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("check_thread_by_id",
		"SELECT id, forum, slug FROM threads WHERE id = $1")
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("check_vote",
		"SELECT voice FROM votes WHERE thread = $1 AND nickname = $2",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("insert_vote",
		"INSERT INTO votes (voice, nickname, thread) "+
			"VALUES ($1, $2, $3) ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("update_vote",
		"UPDATE votes SET voice = $1 "+
			"WHERE nickname = $2 AND thread = $3",
	)
	if err != nil {
		return err
	}

	return nil
}
