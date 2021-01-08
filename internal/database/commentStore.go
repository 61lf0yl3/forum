package database

import (
	"fmt"

	"github.com/astgot/forum/internal/model"
)

// AddComment ...
func (d *Database) AddComment(c *model.Comments) bool {
	stmnt, err := d.db.Prepare("INSERT INTO Comments (postID, author, content, creationDate) VALUES (?, ?, ?, ?)")
	defer stmnt.Close()
	_, err = stmnt.Exec(c.PostID, c.Author, c.Content, c.CreationDate)
	if err != nil {
		return false
	}
	return true
}

// GetCommentByID ...
func (d *Database) GetCommentByID(id int64) (*model.Comments, error) {
	comment := model.NewComment()
	if err := d.db.QueryRow("SELECT postID FROM Comments WHERE ID = ?", id).
		Scan(&comment.PostID); err != nil {
		return nil, err
	}
	return comment, nil
}

// GetCommentsOfPost ...
func (d *Database) GetCommentsOfPost(pid int64) ([]*model.Comments, error) {
	var comments []*model.Comments
	query, err := d.db.Query("SELECT * FROM Comments WHERE postID = ?", pid)
	if err != nil {
		return nil, err
	}
	for query.Next() {
		c := model.NewComment()
		if err = query.Scan(&c.CommentID, &c.PostID, &c.Author, &c.Content, &c.CreationDate); err != nil {
			fmt.Println(err.Error(), "GetCommentsOfPost error")
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}
