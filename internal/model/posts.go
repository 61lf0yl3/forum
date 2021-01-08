package model

// Post ...
type Post struct {
	ID           int64 `db:"postID"`
	UserID       int64 `db:"userID"`
	Title        string
	Author       string
	Content      string
	CreationDate string `db:"creationDate"`
}

// NewPost ...
func NewPost() *Post {
	return &Post{}
}
