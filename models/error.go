package models

import (
	"errors"
)

type Error struct {
	Message string `json:"message,omitempty"`
}

var (
	UserExists      = errors.New("user already exists")
	UserConflict = errors.New("user conflict")
	ForumExists     = errors.New("forum exists")
	ForumNotExists = errors.New("forum does not exists")
	ThreadExists    = errors.New("thread already exists")
	UserNotExists   = errors.New("user does not exist")
	ThreadNotExists = errors.New("no such thread")
	PostNotExists   = errors.New("post does not exist")
	ParentNotExists   = errors.New("patent does not exist")
	SameVote = errors.New("this vote already exists")
)