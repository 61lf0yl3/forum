package model

// PostMap ...
type PostMap struct {
	ThreadID int64 `db:"threadID"`
	PostID   int64 `db:"postID"`
}

// NewPostMap ..
func NewPostMap() *PostMap {
	return &PostMap{}
}
