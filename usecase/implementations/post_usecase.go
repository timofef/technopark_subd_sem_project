package implementations

import (
	"github.com/timofef/technopark_subd_sem_project/models"
	repository "github.com/timofef/technopark_subd_sem_project/repository/interfaces"
	usecase "github.com/timofef/technopark_subd_sem_project/usecase/interfaces"
	"strings"
)

type PostUsecase struct {
	forumRepo  repository.ForumRepository
	userRepo   repository.UserRepository
	threadRepo repository.ThreadRepository
	postRepo   repository.PostRepository
}

func NewPostUsecase(forumR repository.ForumRepository, userR repository.UserRepository,
	threadR repository.ThreadRepository, postR repository.PostRepository) usecase.PostUsecase {
	return &PostUsecase{forumRepo: forumR, userRepo: userR, threadRepo: threadR, postRepo: postR}
}

func (pu *PostUsecase) GetPostDetails(id *string, related []byte) (*models.PostFull, error) {
	postFull := models.PostFull{}

	post, err := pu.postRepo.GetPostById(id)
	if err != nil {
		return nil, err
	}
	postFull.Post = post

	if related != nil {
		relatedArgs := strings.Split(string(related), ",")

		for _, value := range relatedArgs {
			switch value {
			case "user":
				user, _ := pu.userRepo.GetUserByNickname(post.Author)
				postFull.Author = user
			case "thread":
				thread, _ := pu.threadRepo.GetThreadById(post.Thread)
				postFull.Thread = thread
			case "forum":
				forum, _ := pu.forumRepo.GetDetailsBySlug(post.Forum)
				postFull.Forum = forum
			}
		}
	}

	return &postFull, nil
}

func (pu *PostUsecase) EditPost(id *string, update *models.PostUpdate) (*models.Post, error) {
	post, err := pu.postRepo.EditPost(id, update)

	return post, err
}
