package interfaces

type ThreadUsecase interface {
	CreatePost()
	GetThread()
	UpdateThread()
	GetPosts()
	VoteForThread()
}