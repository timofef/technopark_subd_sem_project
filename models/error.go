package models

import (
	"errors"
)

type Error struct {
	Message string `json:"message,omitempty"`
}

var (
	UserExists      = errors.New("user already exist")
	UserNotUpdated = errors.New("can't update user")
	ForumExists     = errors.New("forum exists")
	ForumNotExists = errors.New("forum does not exists")
	ThreadExist    = errors.New("thread already exists")
	UserNotExists   = errors.New("user does not exist")
	ThreadNotExists = errors.New("no such thread")
	SqlError       = errors.New("sql error")
	PostNotExist   = errors.New("post not exist")
)