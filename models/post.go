package models

import "github.com/go-openapi/strfmt"

//easyjson:json
type Post struct {
	Author   string          `json:"author"`
	Created  strfmt.DateTime `json:"created,omitempty"`
	Forum    string          `json:"forum,omitempty"`
	ID       int64           `json:"id,omitempty"`
	IsEdited bool            `json:"isEdited,omitempty"`
	Message  string          `json:"message"`
	Parent   int64           `json:"parent,omitempty"`
	Thread   int32           `json:"thread,omitempty"`
}

//easyjson:json
type Posts []*Post