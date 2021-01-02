package implementations

import (
	"bytes"
	"database/sql"
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
	newRepo := &ThreadRepo{db: pool}

	if err := newRepo.PrepareStatements(); err != nil {
		fmt.Println(err)
	}

	return newRepo
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

	checkUser := tx.QueryRow("check_user_by_nickname", thread.Author)
	if err = checkUser.Scan(&thread.Author); err != nil {
		return nil, models.ForumNotExists
	}

	checkForum := tx.QueryRow("check_forum_by_slug", thread.Forum)
	if err = checkForum.Scan(&thread.Forum); err != nil {
		return nil, models.ForumNotExists
	}

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

func (t *ThreadRepo) GetThreadById(id int32) (*models.Thread, error) {
	tx, err := t.db.Begin()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	rows := tx.QueryRow("get_thread_by_id", id)

	var thread = models.Thread{}
	slug := sql.NullString{}
	err = rows.Scan(&thread.ID,
		&thread.Author,
		&thread.Created,
		&thread.Forum,
		&thread.Message,
		&slug,
		&thread.Title,
		&thread.Votes,
	)
	if slug.Valid {
		thread.Slug = slug.String
	}

	if err != nil {
		//fmt.Println(err)
		return nil, models.ThreadNotExists
	}

	return &thread, nil
}

func (t *ThreadRepo) GetThreadBySlugOrId(slugOrId *interface{}) (*models.Thread, error) {
	tx, err := t.db.Begin()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	var thread = models.Thread{}
	var slug = sql.NullString{}
	id, err := strconv.Atoi((*slugOrId).(string))
	if err != nil {
		if err = tx.QueryRow("get_thread_by_slug", slugOrId).
			Scan(&thread.ID,
				&thread.Author,
				&thread.Created,
				&thread.Forum,
				&thread.Message,
				&slug,
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
				&slug,
				&thread.Title,
				&thread.Votes,
			); err != nil {
			return nil, models.ThreadNotExists
		}
	}
	if slug.Valid {
		thread.Slug = slug.String
	}

	return &thread, nil
}

func CheckThreadBySlugOrId(tx *pgx.Tx, slugOrId *interface{}) (*models.Thread, error) {
	thread := models.Thread{}
	id, err := strconv.Atoi((*slugOrId).(string))
	if err != nil {
		if err = tx.QueryRow("check_thread_by_slug", slugOrId).
			Scan(&thread.ID, &thread.Forum); err != nil {
			return nil, models.ThreadNotExists
		}
	} else {
		if err = tx.QueryRow("check_thread_by_id", id).
			Scan(&thread.ID, &thread.Forum); err != nil {
			return nil, models.ThreadNotExists
		}
	}

	return &thread, nil
}

func (t *ThreadRepo) VoteForThread(slugOrId *interface{}, voice *models.Vote) (*models.Thread, error) {
	tx, err := t.db.Begin()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	thread, err := CheckThreadBySlugOrId(tx, slugOrId)
	if err != nil {
		return nil, err
	}

	checkUser := tx.QueryRow("check_user_by_nickname", voice.Nickname)
	if err = checkUser.Scan(&voice.Nickname); err != nil {
		return nil, models.UserNotExists
	}

	rows := tx.QueryRow("get_thread_by_id", thread.ID)

	slug := sql.NullString{}
	err = rows.Scan(&thread.ID,
		&thread.Author,
		&thread.Created,
		&thread.Forum,
		&thread.Message,
		&slug,
		&thread.Title,
		&thread.Votes,
	)
	if slug.Valid {
		thread.Slug = slug.String
	}

	if err != nil {
		fmt.Println(err)
		return nil, models.ThreadNotExists
	}

	var alreadyVoted int32 = 0

	err = tx.QueryRow("check_vote",
		voice.Nickname,
		thread.ID,
	).Scan(&alreadyVoted)

	if alreadyVoted == 0 {
		_, err = tx.Exec("insert_vote",
			voice.Voice,
			voice.Nickname,
			thread.ID,
		)
		//fmt.Println("err insert   ", err)
		thread.Votes += voice.Voice
	} else {
		if alreadyVoted != voice.Voice {
			_, err = tx.Exec("update_vote",
				voice.Voice,
				voice.Nickname,
				thread.ID,
			)
			//fmt.Println("err update   ", err)
			thread.Votes += 2 * voice.Voice
		}
	}

	return thread, nil
}

func (t *ThreadRepo) GetThreadPostsFlat(tx *pgx.Tx, thread int32, limit, since, desc []byte) (*models.Posts, error) {
	trueBytes := []byte("true")
	var rows *pgx.Rows
	var err error

	switch true {
	case since == nil && limit == nil && !bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_flat", thread)
	case since == nil && limit == nil && bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_flat_desc", thread)
	case since == nil && limit != nil && !bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_flat_limit", thread, limit)
	case since == nil && limit != nil && bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_flat_limit_desc", thread, limit)
	case since != nil && limit == nil && !bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_flat_since", thread, since)
	case since != nil && limit == nil && bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_flat_since_desc", thread, since)
	case since != nil && limit != nil && !bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_flat_since_limit", thread, since, limit)
	case since != nil && limit != nil && bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_flat_since_limit_desc", thread, since, limit)
	}

	if err != nil {
		fmt.Println("getThreadPostsFlat  ", err)
		return &models.Posts{}, err
	}

	posts := models.Posts{}

	for rows.Next() {
		post := models.Post{}

		if err = rows.Scan(
			&post.ID,
			&post.Author,
			&post.Created,
			&post.Forum,
			&post.IsEdited,
			&post.Message,
			&post.Parent,
			&post.Thread);
			err != nil {
			fmt.Println("getThreadPostsFlat  Scan rows  ", err)
			return &models.Posts{}, err
		}
		posts = append(posts, &post)
	}

	return &posts, nil
}

func (t *ThreadRepo) GetThreadPostsTree(tx *pgx.Tx, thread int32, limit, since, desc []byte) (*models.Posts, error) {
	trueBytes := []byte("true")
	var rows *pgx.Rows
	var err error

	switch true {
	case since == nil && limit == nil && !bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_tree", thread)
	case since == nil && limit == nil && bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_tree_desc", thread)
	case since == nil && limit != nil && !bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_tree_limit", thread, limit)
	case since == nil && limit != nil && bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_tree_limit_desc", thread, limit)
	case since != nil && limit == nil && !bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_tree_since", thread, since)
	case since != nil && limit == nil && bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_tree_since_desc", thread, since)
	case since != nil && limit != nil && !bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_tree_since_limit", thread, since, limit)
	case since != nil && limit != nil && bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_tree_since_limit_desc", thread, since, limit)
	}

	if err != nil {
		fmt.Println("getThreadPostsTree  ", err, string(since), string(limit), string(desc))
		return &models.Posts{}, err
	}

	posts := models.Posts{}

	for rows.Next() {
		post := models.Post{}

		if err = rows.Scan(
			&post.ID,
			&post.Author,
			&post.Created,
			&post.Forum,
			&post.IsEdited,
			&post.Message,
			&post.Parent,
			&post.Thread);
			err != nil {
			fmt.Println("getThreadPostsTree  Scan rows  ", err)
			return &models.Posts{}, err
		}
		posts = append(posts, &post)
	}

	return &posts, nil
}

func (t *ThreadRepo) GetThreadPostsParentTree(tx *pgx.Tx, thread int32, limit, since, desc []byte) (*models.Posts, error) {
	trueBytes := []byte("true")
	var rows *pgx.Rows
	var err error

	switch true {
	case since == nil && limit == nil && !bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_parenttree", thread)
	case since == nil && limit == nil && bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_parenttree_desc", thread)
	case since == nil && limit != nil && !bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_parenttree_limit", thread, limit)
	case since == nil && limit != nil && bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_parenttree_limit_desc", thread, limit)
	case since != nil && limit == nil && !bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_parenttree_since", thread, since)
	case since != nil && limit == nil && bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_parenttree_since_desc", thread, since)
	case since != nil && limit != nil && !bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_parenttree_since_limit", thread, since, limit)
	case since != nil && limit != nil && bytes.Equal(desc, trueBytes):
		rows, err = tx.Query("get_posts_parenttree_since_limit_desc", thread, since, limit)
	}

	if err != nil {
		fmt.Println("getThreadPostsParentTree  ", err, string(since), string(limit), string(desc))
		return &models.Posts{}, err
	}

	posts := models.Posts{}

	for rows.Next() {
		post := models.Post{}

		if err = rows.Scan(
			&post.ID,
			&post.Author,
			&post.Created,
			&post.Forum,
			&post.IsEdited,
			&post.Message,
			&post.Parent,
			&post.Thread);
			err != nil {
			fmt.Println("getThreadPostsParentTree  Scan rows  ", err)
			return &models.Posts{}, err
		}
		posts = append(posts, &post)
	}

	return &posts, nil
}

func (t *ThreadRepo) GetThreadPosts(slugOrId *interface{}, limit, since, sort, desc []byte) (*models.Posts, error) {
	tx, err := t.db.Begin()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	thread, err := CheckThreadBySlugOrId(tx, slugOrId)
	if err != nil {
		return nil, err
	}

	switch true {
	case bytes.Equal([]byte("tree"), sort):
		posts, err := t.GetThreadPostsTree(tx, thread.ID, limit, since, desc)
		return posts, err
	case bytes.Equal([]byte("parent_tree"), sort):
		posts, err := t.GetThreadPostsParentTree(tx, thread.ID, limit, since, desc)
		return posts, err
	default:
		posts, err := t.GetThreadPostsFlat(tx, thread.ID, limit, since, desc)
		return posts, err
	}
}

func (t *ThreadRepo) UpdateThreadById(slugOrId *interface{}, threadUpdate *models.ThreadUpdate) (*models.Thread, error) {
	tx, err := t.db.Begin()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	thread, err := CheckThreadBySlugOrId(tx, slugOrId)
	if err != nil {
		return nil, err
	}

	rows := tx.QueryRow("update_thread",
		sql.NullString{String: threadUpdate.Message, Valid: threadUpdate.Message != ""},
		sql.NullString{String: threadUpdate.Title, Valid: threadUpdate.Title != ""},
		thread.ID)

	if err = rows.Scan(&thread.ID,
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

	if threadUpdate.Message != "" {
		thread.Message = threadUpdate.Message
	}
	if threadUpdate.Title != "" {
		thread.Title = threadUpdate.Title
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
		"SELECT id, forum FROM threads WHERE slug = $1")
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("check_thread_by_id",
		"SELECT id, forum FROM threads WHERE id = $1")
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("check_vote",
		"SELECT voice FROM votes WHERE thread = $2 AND nickname = $1",
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

	_, err = t.db.Prepare("update_thread",
		"UPDATE threads SET "+
			"message = COALESCE($1, threads.message), "+
			"title = COALESCE($2, threads.title) "+
			"WHERE id = $3 RETURNING threads.id, threads.author, threads.created, threads.forum, "+
			"threads.message, threads.slug, threads.title, threads.votes",
	)
	if err != nil {
		return err
	}

	// FLAT

	_, err = t.db.Prepare("get_posts_flat",
		"SELECT id, author, created, forum, is_edited, message, parent, thread "+
			"FROM posts "+
			"WHERE thread = $1::INT "+
			"ORDER BY created, id ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_flat_desc",
		"SELECT id, author, created, forum, is_edited, message, parent, thread "+
			"FROM posts "+
			"WHERE thread = $1::INT "+
			"ORDER BY created DESC, id DESC",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_flat_limit",
		"SELECT id, author, created, forum, is_edited, message, parent, thread "+
			"FROM posts "+
			"WHERE thread = $1::INT "+
			"ORDER BY created, id "+
			"LIMIT $2::TEXT::INT",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_flat_limit_desc",
		"SELECT id, author, created, forum, is_edited, message, parent, thread "+
			"FROM posts "+
			"WHERE thread = $1::INT "+
			"ORDER BY created DESC, id DESC "+
			"LIMIT $2::TEXT::INT",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_flat_since",
		"SELECT id, author, created, forum, is_edited, message, parent, thread "+
			"FROM posts "+
			"WHERE thread = $1::INT AND id > $2::TEXT::INT "+
			"ORDER BY created, id ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_flat_since_desc",
		"SELECT id, author, created, forum, is_edited, message, parent, thread "+
			"FROM posts "+
			"WHERE thread = $1::INT AND id < $2::TEXT::INT "+
			"ORDER BY created DESC, id DESC ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_flat_since_limit",
		"SELECT id, author, created, forum, is_edited, message, parent, thread "+
			"FROM posts "+
			"WHERE thread = $1::INT AND id > $2::TEXT::INT "+
			"ORDER BY created, id "+
			"LIMIT $3::TEXT::INT",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_flat_since_limit_desc",
		"SELECT id, author, created, forum, is_edited, message, parent, thread "+
			"FROM posts "+
			"WHERE thread = $1::INT AND id < $2::TEXT::INT "+
			"ORDER BY created DESC, id DESC "+
			"LIMIT $3::TEXT::INT",
	)
	if err != nil {
		return err
	}

	// TREE

	_, err = t.db.Prepare("get_posts_tree",
		"SELECT id, author, created, forum, is_edited, message, parent, thread "+
			"FROM posts "+
			"WHERE thread = $1::INT "+
			"ORDER BY path ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_tree_desc",
		"SELECT id, author, created, forum, is_edited, message, parent, thread "+
			"FROM posts "+
			"WHERE thread = $1::INT "+
			"ORDER BY path DESC",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_tree_limit",
		"SELECT id, author, created, forum, is_edited, message, parent, thread "+
			"FROM posts "+
			"WHERE thread = $1::INT "+
			"ORDER BY path "+
			"LIMIT $2::TEXT::INT",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_tree_limit_desc",
		"SELECT id, author, created, forum, is_edited, message, parent, thread "+
			"FROM posts "+
			"WHERE thread = $1::INT "+
			"ORDER BY path DESC "+
			"LIMIT $2::TEXT::INT",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_tree_since",
		"SELECT id, author, created, forum, is_edited, message, parent, thread "+
			"FROM posts "+
			"WHERE thread = $1::INT AND path > (SELECT path FROM posts WHERE id = $2::TEXT::INT) "+
			"ORDER BY path ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_tree_since_desc",
		"SELECT id, author, created, forum, is_edited, message, parent, thread "+
			"FROM posts "+
			"WHERE thread = $1::INT  AND path < (SELECT path FROM posts WHERE id = $2::TEXT::INT) "+
			"ORDER BY path DESC",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_tree_since_limit",
		"SELECT id, author, created, forum, is_edited, message, parent, thread "+
			"FROM posts "+
			"WHERE thread = $1::INT AND path > (SELECT path FROM posts WHERE id = $2::TEXT::INT) "+
			"ORDER BY path "+
			"LIMIT $3::TEXT::INT",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_tree_since_limit_desc",
		"SELECT id, author, created, forum, is_edited, message, parent, thread "+
			"FROM posts "+
			"WHERE thread = $1::INT  AND path < (SELECT path FROM posts WHERE id = $2::TEXT::INT) "+
			"ORDER BY path DESC "+
			"LIMIT $3::TEXT::INT ",
	)
	if err != nil {
		return err
	}

	// PARENT TREE

	_, err = t.db.Prepare("get_posts_parenttree",
		"SELECT posts.id, posts.author, posts.created, posts.forum, posts.is_edited, posts.message, posts.parent, posts.thread "+
			"FROM (SELECT * FROM posts r WHERE r.parent = 0 AND r.thread = $1::INT "+
			"ORDER BY r.path) AS root "+
			"JOIN posts ON root.path[1] = posts.path[1] "+
			"ORDER BY posts.path[1], posts.path ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_parenttree_desc",
		"SELECT posts.id, posts.author, posts.created, posts.forum, posts.is_edited, posts.message, posts.parent, posts.thread "+
			"FROM (SELECT * FROM posts r WHERE r.parent = 0 AND r.thread = $1::INT "+
			"ORDER BY r.path DESC) AS root "+
			"JOIN posts ON root.path[1] = posts.path[1] "+
			"ORDER BY posts.path[1] DESC, posts.path  ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_parenttree_limit",
		"SELECT posts.id, posts.author, posts.created, posts.forum, posts.is_edited, posts.message, posts.parent, posts.thread "+
			"FROM (SELECT * FROM posts r WHERE r.parent = 0 AND r.thread = $1::INT "+
			"ORDER BY r.path LIMIT $2::TEXT::INT) AS root "+
			"JOIN posts ON root.path[1] = posts.path[1] "+
			"ORDER BY posts.path[1], posts.path  ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_parenttree_limit_desc",
		"SELECT posts.id, posts.author, posts.created, posts.forum, posts.is_edited, posts.message, posts.parent, posts.thread "+
			"FROM (SELECT * FROM posts r WHERE r.parent = 0 AND r.thread = $1::INT"+
			" ORDER BY r.path DESC LIMIT $2::TEXT::INT) AS root "+
			"JOIN posts ON root.path[1] = posts.path[1] "+
			"ORDER BY posts.path[1] DESC, posts.path  ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_parenttree_since",
		"SELECT posts.id, posts.author, posts.created, posts.forum, posts.is_edited, posts.message, posts.parent, posts.thread "+
			"FROM (SELECT * FROM posts r WHERE r.parent = 0 AND r.thread = $1::INT "+
			"AND r.path[1] > (SELECT path[1] FROM posts WHERE id = $2::TEXT::INT) ORDER BY r.path) AS root "+
			"JOIN posts ON root.path[1] = posts.path[1] "+
			"ORDER BY posts.path[1], posts.path  ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_parenttree_since_desc",
		"SELECT posts.id, posts.author, posts.created, posts.forum, posts.is_edited, posts.message, posts.parent, posts.thread "+
			"FROM (SELECT * FROM posts r WHERE r.parent = 0 AND r.thread = $1::INT "+
			"AND r.path[1] < (SELECT path[1] FROM posts WHERE id = $2::TEXT::INT) ORDER BY r.path DESC ) AS root "+
			"JOIN posts ON root.path[1] = posts.path[1] "+
			"ORDER BY posts.path[1] DESC, posts.path  ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_parenttree_since_limit",
		"SELECT posts.id, posts.author, posts.created, posts.forum, posts.is_edited, posts.message, posts.parent, posts.thread "+
			"FROM (SELECT * FROM posts r WHERE r.parent = 0 AND r.thread = $1::INT "+
			"AND r.path[1] > (SELECT path[1] FROM posts WHERE id = $2::TEXT::INT) "+
			"ORDER BY r.path  LIMIT $3::TEXT::INT) AS root "+
			"JOIN posts ON root.path[1] = posts.path[1] "+
			"ORDER BY posts.path[1], posts.path  ",
	)
	if err != nil {
		return err
	}

	_, err = t.db.Prepare("get_posts_parenttree_since_limit_desc",
		"SELECT posts.id, posts.author, posts.created, posts.forum, posts.is_edited, posts.message, posts.parent, posts.thread "+
			"FROM (SELECT * FROM posts r WHERE r.parent = 0 AND r.thread = $1::INT "+
			"AND r.path[1] < (SELECT path[1] FROM posts WHERE id = $2::TEXT::INT) "+
			"ORDER BY r.path DESC LIMIT $3::TEXT::INT) AS root "+
			"JOIN posts ON root.path[1] = posts.path[1] "+
			"ORDER BY posts.path[1] DESC, posts.path  ",
	)
	if err != nil {
		return err
	}

	return nil
}
