package models

import (
	"github.com/go-openapi/strfmt"
)

type User struct {
	About    string       `json:"about,omitempty"`
	Email    strfmt.Email `json:"email"`
	Fullname string       `json:"fullname"`
	Nickname string       `json:"nickname,omitempty"`
}

type Users []*User
