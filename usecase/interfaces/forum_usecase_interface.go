package interfaces

import "github.com/timofef/technopark_subd_sem_project/models"

type ForumUsecase interface {
	CreateForum(forum *models.Forum) (*models.Forum, error)
	CreateThread(thread *models.Thread) (*models.Thread, error)
	GetForumDetails(slug string) (*models.Forum, error)
	GetForumThreads(slug string, since []byte, desc []byte, limit []byte) (*models.Threads, error)
	GetForumUsers()
}