package models

import "github.com/go-openapi/strfmt"

type UserUpdate struct {
	About    string       `json:"about,omitempty"`
	Email    strfmt.Email `json:"email,omitempty"`
	Fullname string       `json:"fullname,omitempty"`
}
