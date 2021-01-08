package database

import (
	"fmt"

	"github.com/astgot/forum/internal/model"
)

// InsertPostInfo ...
func (d *Database) InsertPostInfo(p *model.Post) (int64, error) {
	if err := d.Open(); err != nil {
		return -1, err
	}

	stmnt, err := d.db.Prepare("INSERT INTO Posts (user_id, author, title, content, creationDate) VALUES (?, ?, ?, ?, ?)")
	defer stmnt.Close()
	if err != nil {
		return -1, err
	}
	res, err := stmnt.Exec(p.UserID, p.Author, p.Title, p.Content, p.CreationDate)
	if err != nil {
		return -1, err
	}
	// to get PostID of new post for error handling in DB
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

// GetPosts ...
func (d *Database) GetPosts() ([]*model.Post, error) {
	var posts []*model.Post
	if err := d.Open(); err != nil {
		return nil, err
	}
	// show posts in reverse order (in the beginning fresh posts)
	res, err := d.db.Query("SELECT * FROM Posts ORDER BY post_id DESC") // DESC - in reverse order
	if err != nil {
		return nil, err
	}
	defer res.Close()
	for res.Next() {
		post := model.NewPost()
		if err := res.Scan(&post.ID, &post.UserID, &post.Author, &post.Title, &post.Content, &post.CreationDate); err != nil {
			fmt.Println(err.Error(), "GetPosts() postsStore.go")
			return nil, err

		}
		posts = append(posts, post)

	}
	return posts, nil
}

// GetPostByPID ...
func (d *Database) GetPostByPID(pid int64) (*model.Post, error) {
	post := model.NewPost()
	if err := d.db.QueryRow("SELECT author, title, content, creationDate FROM Posts WHERE post_id = ?", pid).Scan(
		&post.Author,
		&post.Title,
		&post.Content,
		&post.CreationDate,
	); err != nil {
		return nil, err
	}
	post.ID = pid
	return post, nil
}

// GetPostsByUID ...
func (d *Database) GetPostsByUID(uid int64) ([]*model.Post, error) {
	var posts []*model.Post
	if err := d.Open(); err != nil {
		return nil, err
	}
	query, err := d.db.Query("SELECT * FROM Posts WHERE user_id=? ORDER BY post_id DESC", uid)
	if err != nil {
		fmt.Println("GetPostsByUID error", err.Error())
		return nil, err
	}
	defer query.Close()
	for query.Next() {
		post := model.NewPost()
		if err := query.Scan(&post.ID, &post.UserID, &post.Author, &post.Title, &post.Content, &post.CreationDate); err != nil {
			fmt.Println(err.Error(), "GetPostsByUID error")
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
