package implementations

import (
	"database/sql"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx"
	"github.com/timofef/technopark_subd_sem_project/models"
	"github.com/timofef/technopark_subd_sem_project/repository/interfaces"
	"time"
)

type PostRepo struct {
	db *pgx.ConnPool
}

func NewPostRepo(pool *pgx.ConnPool) interfaces.PostRepository {
	new := &PostRepo{db: pool}

	if err := new.PrepareStatements(); err != nil {
		fmt.Println(err)
	}

	return new
}

func (p *PostRepo) CreatePosts(slugOrId *interface{}, posts *models.Posts) (*models.Posts, error) {
	tx, err := p.db.Begin()
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

	if len(*posts) == 0 {
		return &models.Posts{}, nil
	}

	creationTime := strfmt.DateTime(time.Now())
	for _, post := range *posts {
		post.Thread = thread.ID
		post.Forum = thread.Forum
		post.Created = creationTime

		err = tx.QueryRow("insert_post",
			post.Author,
			creationTime,
			post.Forum,
			post.Message,
			post.Parent,
			post.Thread,
		).Scan(&post.ID)

		if err != nil {
			//fmt.Println(err)
			if err.Error() == "ERROR: 00404 (SQLSTATE 00404)" {
				return nil, models.ParentNotExists
			} else {
				return nil, models.ThreadNotExists
			}

		}
	}

	return posts, nil
}

func (p *PostRepo) GetPostById(id *string) (*models.Post, error) {
	tx, err := p.db.Begin()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	rows := tx.QueryRow("get_post_by_id", *id)

	var post = models.Post{}
	err = rows.Scan(&post.ID,
		&post.Author,
		&post.Created,
		&post.Forum,
		&post.IsEdited,
		&post.Message,
		&post.Parent,
		&post.Thread,
	)

	if err != nil {
		//fmt.Println(err)
		return nil, models.PostNotExists
	}

	return &post, nil
}

func (p *PostRepo) EditPost(id *string, update *models.PostUpdate) (*models.Post, error) {
	tx, err := p.db.Begin()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	post := models.Post{}

	if err := tx.QueryRow("update_post",
		sql.NullString{String: update.Message, Valid: update.Message != ""},
		id).
		Scan(&post.ID,
			&post.Author,
			&post.Created,
			&post.Forum,
			&post.IsEdited,
			&post.Message,
			&post.Parent,
			&post.Thread); err != nil {
		//fmt.Println(err)
		return nil, models.PostNotExists
	}

	return &post, nil
}

func (p *PostRepo) PrepareStatements() error {
	_, err := p.db.Prepare("insert_post",
		"INSERT INTO posts (author, created, forum, message, parent, thread) "+
			"VALUES ($1, $2, $3, $4, $5, $6) "+
			"RETURNING id ",
	)
	if err != nil {
		return err
	}

	_, err = p.db.Prepare("get_post_by_id",
		"SELECT id, author, created, forum, is_edited, message, parent, thread "+
			"FROM posts "+
			"WHERE id = $1 ",
	)
	if err != nil {
		return err
	}

	_, err = p.db.Prepare("update_post",
		"UPDATE posts SET "+
			"message = COALESCE($1, message), "+
			"is_edited = (CASE WHEN $1 IS NULL OR $1 = message THEN FALSE ELSE TRUE END) "+
			"WHERE id = $2 "+
			"RETURNING id, author, created, forum, is_edited, message, parent, thread ",
	)
	if err != nil {
		return err
	}

	return nil
}
