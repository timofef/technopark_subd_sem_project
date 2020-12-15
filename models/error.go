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
	UserConflict = errors.New("user conflict")
	ForumExists     = errors.New("forum exists")
	ForumNotExists = errors.New("forum does not exists")
	ThreadExists    = errors.New("thread already exists")
	UserNotExists   = errors.New("user does not exist")
	ThreadNotExists = errors.New("no such thread")
	SqlError       = errors.New("sql error")
	PostNotExists   = errors.New("post not exist")
)