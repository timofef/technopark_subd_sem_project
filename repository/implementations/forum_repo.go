package implementations

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/timofef/technopark_subd_sem_project/models"
	"github.com/timofef/technopark_subd_sem_project/repository/interfaces"
)

type ForumRepo struct {
	db *pgx.ConnPool
}

func NewForumRepo(pool *pgx.ConnPool) interfaces.ForumRepository {
	new := &ForumRepo{db: pool}

	if err := new.PrepareStatements(); err != nil {
		fmt.Println(err)
	}

	return new
}

func (f *ForumRepo) CreateForum(forum *models.Forum) (*models.Forum, error) {
	tx, err := f.db.Begin()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	result, err := tx.Exec("insert_forum", forum.Slug, forum.Title, forum.User)
	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		existingForum := models.Forum{Slug: forum.Slug}

		rows := tx.QueryRow("get_forum_by_slug", forum.Slug)

		_ = rows.Scan(&existingForum.Title,
			&existingForum.User,
			&existingForum.Slug,
			&existingForum.Posts,
			&existingForum.Threads)

		return &existingForum, models.ForumExists
	}

	return forum, nil
}

func (f *ForumRepo) GetDetailsBySlug(slug string) (*models.Forum, error) {
	tx, err := f.db.Begin()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	rows := tx.QueryRow("get_forum_by_slug", slug)

	var forum = models.Forum{}
	err = rows.Scan(&forum.Title,
		&forum.User,
		&forum.Slug,
		&forum.Posts,
		&forum.Threads)

	if err != nil {
		return nil, models.ForumNotExists
	}

	return &forum, nil
}



func (f *ForumRepo) GetUsersBySlug() {
	panic("implement me")
}

func (f *ForumRepo) GetTreadsBySlug() {
	panic("implement me")
}

func (f *ForumRepo) PrepareStatements() error {
	_, err := f.db.Prepare("insert_forum",
		"INSERT INTO forums (slug, title, owner) "+
			"VALUES ($1, $2, $3) "+
			"ON CONFLICT DO NOTHING ",
	)
	if err != nil {
		return err
	}

	_, err = f.db.Prepare("get_forum_by_slug",
		"SELECT forums.title, forums.owner, forums.slug, forums.posts, forums.threads "+
			"FROM forums "+
			"WHERE slug = $1 ",
	)
	if err != nil {
		return err
	}



	return nil
}
