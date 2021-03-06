package implementations

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/timofef/technopark_subd_sem_project/models"
	"github.com/timofef/technopark_subd_sem_project/repository/interfaces"
)

type ForumRepo struct {
	db *pgx.ConnPool
}

func NewForumRepo(pool *pgx.ConnPool) interfaces.ForumRepository {
	newRepo := &ForumRepo{db: pool}

	if err := newRepo.PrepareStatements(); err != nil {
		fmt.Println(err)
	}

	return newRepo
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

	checkUser := tx.QueryRow("check_user_by_nickname", forum.User)
	if err = checkUser.Scan(&forum.User); err != nil {
		return nil, models.UserNotExists
	}

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

func (f *ForumRepo) GetDetailsBySlug(slug *string) (*models.Forum, error) {
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

func (f *ForumRepo) GetThreads(slug *string, since, desc, limit []byte) (*models.Threads, error) {
	tx, err := f.db.Begin()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	checkForum := tx.QueryRow("check_forum_by_slug", slug)
	if err = checkForum.Scan(slug); err != nil {
		return nil, models.ForumNotExists
	}

	var rows *pgx.Rows

	if since == nil {
		if bytes.Equal([]byte("true"), desc) {
			rows, err = tx.Query("get_threads_limit_desc", slug, limit)
		} else {
			rows, err = tx.Query("get_threads_limit", slug, limit)
		}
	} else {
		if bytes.Equal([]byte("true"), desc) {
			rows, err = tx.Query("get_threads_limit_since_desc", slug, since, limit)
		} else {
			rows, err = tx.Query("get_threads_limit_since", slug, since, limit)
		}
	}

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	threads := models.Threads{}
	for rows.Next() {
		curThread := models.Thread{}
		slug := sql.NullString{}
		err = rows.Scan(&curThread.ID,
			&curThread.Author,
			&curThread.Created,
			&curThread.Forum,
			&curThread.Message,
			&slug,
			&curThread.Title,
			&curThread.Votes)

		if slug.Valid {
			curThread.Slug = slug.String
		}

		threads = append(threads, curThread)
	}
	rows.Close()

	return &threads, nil
}

func (f *ForumRepo) GetUsersBySlug(slug *string, since, desc, limit []byte) (*models.Users, error) {
	tx, err := f.db.Begin()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	checkForum := tx.QueryRow("check_forum_by_slug", slug)
	if err = checkForum.Scan(slug); err != nil {
		return nil, models.ForumNotExists
	}

	var rows *pgx.Rows

	if since == nil {
		if bytes.Equal([]byte("true"), desc) {
			rows, err = tx.Query("get_users_limit_desc", slug, limit)
		} else {
			rows, err = tx.Query("get_users_limit", slug, limit)
		}
	} else {
		if bytes.Equal([]byte("true"), desc) {
			rows, err = tx.Query("get_users_limit_since_desc", slug, string(since), limit)
		} else {
			rows, err = tx.Query("get_users_limit_since", slug, string(since), limit)
		}
	}

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	users := models.Users{}
	for rows.Next() {
		curUser := models.User{}
		_ = rows.Scan(&curUser.Email,
			&curUser.Fullname,
			&curUser.Nickname,
			&curUser.About)

		users = append(users, &curUser)
	}
	rows.Close()

	return &users, nil
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

	_, err = f.db.Prepare("check_forum_by_slug",
		"SELECT forums.slug "+
			"FROM forums "+
			"WHERE slug = $1 ",
	)
	if err != nil {
		return err
	}

	// GET THREADs
	_, err = f.db.Prepare("get_threads_limit",
		"SELECT id, author, created, forum, message, slug, title, votes "+
			"FROM threads "+
			"WHERE forum = $1::TEXT "+
			"ORDER BY created "+
			"LIMIT $2::TEXT::INT",
	)
	if err != nil {
		return err
	}

	_, err = f.db.Prepare("get_threads_limit_desc",
		"SELECT id, author, created, forum, message, slug, title, votes "+
			"FROM threads "+
			"WHERE forum = $1::TEXT "+
			"ORDER BY created DESC "+
			"LIMIT $2::TEXT::INT",
	)
	if err != nil {
		return err
	}

	_, err = f.db.Prepare("get_threads_limit_since",
		"SELECT id, author, created, forum, message, slug, title, votes "+
			"FROM threads "+
			"WHERE forum = $1::TEXT AND created >= $2::TEXT::TIMESTAMPTZ "+
			"ORDER BY created "+
			"LIMIT $3::TEXT::INT",
	)
	if err != nil {
		return err
	}

	_, err = f.db.Prepare("get_threads_limit_since_desc",
		"SELECT id, author, created, forum, message, slug, title, votes "+
			"FROM threads "+
			"WHERE forum = $1::TEXT AND created <= $2::TEXT::TIMESTAMPTZ "+
			"ORDER BY created DESC "+
			"LIMIT $3::TEXT::INT",
	)
	if err != nil {
		return err
	}

	// GET USERS

	_, err = f.db.Prepare("get_users_limit",
		"SELECT users.email, users.fullname, users.nickname, users.about "+
			"FROM forum_users JOIN users ON forum_users.nickname = users.nickname "+
			"WHERE forum_users.forum = $1::TEXT "+
			"ORDER BY users.nickname "+
			"LIMIT $2::TEXT::INT",
	)
	if err != nil {
		return err
	}

	_, err = f.db.Prepare("get_users_limit_desc",
		"SELECT users.email, users.fullname, users.nickname, users.about "+
			"FROM forum_users JOIN users ON forum_users.nickname = users.nickname "+
			"WHERE forum_users.forum = $1::TEXT "+
			"ORDER BY users.nickname DESC "+
			"LIMIT $2::TEXT::INT",
	)
	if err != nil {
		return err
	}

	_, err = f.db.Prepare("get_users_limit_since",
		"SELECT users.email, users.fullname, users.nickname, users.about "+
			"FROM forum_users JOIN users ON forum_users.nickname = users.nickname "+
			"WHERE forum_users.forum = $1::TEXT AND users.nickname > $2 "+
			"ORDER BY users.nickname "+
			"LIMIT $3::TEXT::INT ",
	)
	if err != nil {
		return err
	}

	_, err = f.db.Prepare("get_users_limit_since_desc",
		"SELECT users.email, users.fullname, users.nickname, users.about "+
			"FROM forum_users JOIN users ON forum_users.nickname = users.nickname "+
			"WHERE forum_users.forum = $1::TEXT AND users.nickname < $2 "+
			"ORDER BY users.nickname DESC "+
			"LIMIT $3::TEXT::INT",
	)
	if err != nil {
		return err
	}

	return nil
}
