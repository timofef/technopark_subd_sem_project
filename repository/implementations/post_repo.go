package implementations

import (
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

func (p *PostRepo) CreatePosts(posts *models.Posts, thread *models.Thread) (*models.Posts, error) {
	if len(*posts) == 0 {
		return &models.Posts{}, nil
	}

	tx, err := p.db.Begin()
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

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
			fmt.Println(err)
			return nil, err
		}
	}

	return posts, nil
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

	return nil
}
