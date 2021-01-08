package model

// PostRating ...
type PostRating struct {
	PostID int64 `db:"postID"`
	// UID          int64 `db:"userID"`
	Like         bool
	LikeCount    int64 `db:"likeCount"`
	DislikeCount int64 `db:"dislikeCount"`
}

// CommentRating ...
type CommentRating struct {
	CommentID int64 `db:"commentID"`
	PostID    int64 `db:"postID"`
	// UID          int64 `db:"userID"`
	Like         bool
	LikeCount    int64 `db:"likeCount"`
	DislikeCount int64 `db:"dislikeCount"`
}

// NewPostRating ...
func NewPostRating() *PostRating {
	return &PostRating{}
}

// NewCommentRating ...
func NewCommentRating() *CommentRating {
	return &CommentRating{}
}
