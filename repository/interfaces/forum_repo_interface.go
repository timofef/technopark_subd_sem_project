package interfaces

type ForumRepository interface {
	CreateForum()
	GetDetailsBySlug()
	CreateBranchBySlug()
	GetUsersBySlug()
	GetTreadsBySlug()
	PrepareStatements()
}
