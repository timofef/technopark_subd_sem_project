package implementations

import (
	"github.com/timofef/technopark_subd_sem_project/models"
	repository "github.com/timofef/technopark_subd_sem_project/repository/interfaces"
	usecase "github.com/timofef/technopark_subd_sem_project/usecase/interfaces"
)

type PostUsecase struct {
	forumRepo repository.ForumRepository
	userRepo repository.UserRepository
	threadRepo repository.ThreadRepository
	postRepo repository.PostRepository
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

	return &postFull, nil
}

func (pu *PostUsecase) EditPost(id *string, update *models.PostUpdate) (*models.Post, error) {
	post, err := pu.postRepo.EditPost(id, update)

	return post, err
}