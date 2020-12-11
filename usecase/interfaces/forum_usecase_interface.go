package interfaces

import "github.com/timofef/technopark_subd_sem_project/models"

type ForumUsecase interface {
	CreateForum(forum *models.Forum) (models.Forum, error)
	CreateThread()
	GetForumDetails()
	GetForumThreads()
	GetForumUsers()
}