package interfaces

type ForumUsecase interface {
	CreateForum()
	CreateThread()
	GetForumDetails()
	GetForumThreads()
	GetForumUsers()
}