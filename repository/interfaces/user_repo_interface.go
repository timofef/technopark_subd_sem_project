package interfaces

type UserRepository interface {
	CreateUser()
	GetUserByNickname()
	UpdateUserByNickname()
	PrepareStatements()
}
