package models

import (
	"github.com/go-openapi/strfmt"
)

//easyjson:json
type User struct {
	About    string       `json:"about,omitempty"`
	Email    strfmt.Email `json:"email"`
	Fullname string       `json:"fullname"`
	Nickname string       `json:"nickname,omitempty"`
}

//easyjson:json
type Users []*User
