package implementations

import (
	"github.com/jackc/pgx"
	"github.com/timofef/technopark_subd_sem_project/models"
	"github.com/timofef/technopark_subd_sem_project/repository/interfaces"
)

type ForumRepo struct {
	db *pgx.ConnPool
}

func NewForumRepo(pool *pgx.ConnPool) interfaces.ForumRepository {
	return &ForumRepo{db: pool}
}

func (f *ForumRepo) CreateForum(forum *models.Forum) (models.Forum, error) {
	tx, err := f.db.Begin()
	defer tx.Commit()

	result, err := tx.Exec("insert_forum", forum.Slug, forum.Title, forum.User)
	if err != nil {
		return models.Forum{}, err
	}

	if result.RowsAffected() == 0 {
		existingForum := models.Forum{Slug: forum.Slug}

		rows := tx.QueryRow("get_forum_by_slug", forum.Slug)

		_ = rows.Scan(&existingForum.Title,
			&existingForum.User,
			&existingForum.Slug,
			&existingForum.Posts,
			&existingForum.Threads)

		return existingForum, models.ForumExists
	}

	return *forum, nil
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
