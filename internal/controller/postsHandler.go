package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/astgot/forum/internal/model"
)

// CreatePostHandler ...
func (m *Multiplexer) CreatePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/create" {
			WarnMessage(w, "404 Not Found")
			return
		}
		// Check user authorization
		cookie, err := r.Cookie("authenticated")
		if err != nil {
			WarnMessage(w, "You need to be authorized")
			// http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
		}

		u, err := m.db.GetUserByCookie(cookie.Value)
		if err != nil {
			WarnMessage(w, "Something went wrong")
			fmt.Println("GetUserByCookie error")
			// http.Error(w, "Something went wrong", http.StatusInternalServerError) // Check workflow of DB
			return
		}
		var Create struct {
			Errors   map[string]string
			Username string
		}
		Create.Username = u.Username
		// Gathering post data
		post := model.NewPost()
		thread := model.NewThread()
		if r.Method == "POST" {
			r.ParseForm()
			post.UserID = u.ID
			post.Author = u.Firstname + " " + u.Lastname + " aka " + "\"" + u.Username + "\""
			post.Title = r.PostFormValue("title")
			post.Content = r.PostFormValue("postContent")
			thread.Name = r.PostFormValue("thread")

			Create.Errors = make(map[string]string)
			if post.Title == "" {
				Create.Errors["Title"] = "\"Title\" field is empty"
				tpl.ExecuteTemplate(w, "postCreate.html", Create)
				return
			} else if post.Content == "" {
				Create.Errors["Content"] = "\"Content\" field is empty"
				tpl.ExecuteTemplate(w, "postCreate.html", Create)
				return
			} else if thread.Name == "" {
				Create.Errors["Category"] = "\"Category\" field is empty"
				tpl.ExecuteTemplate(w, "postCreate.html", Create)
				return
			}
			threads := CheckNumberOfThreads(thread.Name)
			post.CreationDate = time.Now().Format("January 2 15:04")
			post.ID, _ = m.db.InsertPostInfo(post)
			if post.ID == -1 {
				fmt.Println("post.ID == -1")
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}
			// If post has several threads, to this post will attach this info
			for _, threadName := range threads {
				m.db.InsertThreadInfo(threadName, post.ID)
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)

		} else if r.Method == "GET" {
			tpl.ExecuteTemplate(w, "postCreate.html", nil)
		}

	}
}

// PostView ... (/post?id=)
//(single post viewing -> to see comments, rate count OR if user is authenticated, he able to add comments and rate post here)
func (m *Multiplexer) PostView() http.HandlerFunc {

	type PostAttr struct {
		Threads  []*model.Thread
		Comments []*model.Comments
		//Likes, Dislikes
	}
	var singlePost struct {
		PostInfo []*PostAttr
		AuthUser *model.Users
		Post     *model.Post
	}
	return func(w http.ResponseWriter, r *http.Request) {

		id, errID := strconv.Atoi(r.URL.Query().Get("id"))
		if errID != nil {
			fmt.Println("errID != nil")
			WarnMessage(w, "Invalid input id")
			// http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		cookie, err := r.Cookie("authenticated")
		if err != nil {
			// If user is guest
			postAttr := &PostAttr{}
			singlePost.Post, err = m.db.GetPostByPID(int64(id))
			if err != nil {
				fmt.Println("Error on PostView() function")
				WarnMessage(w, "The post not found")
				// http.Error(w, "The post not found", http.StatusNotFound)
				return
			}
			postAttr.Comments, err = m.db.GetCommentsOfPost(int64(id))
			commentRate := model.NewCommentRating()
			for i := 0; i < len(postAttr.Comments); i++ {
				commentRate = m.db.GetRateCountOfComment(postAttr.Comments[i].CommentID, int64(id))
				postAttr.Comments[i].LikeCnt = commentRate.LikeCount
				postAttr.Comments[i].DislikeCnt = commentRate.DislikeCount
			}
			postAttr.Threads, _ = m.db.GetThreadOfPost(int64(id))
			singlePost.PostInfo = append(singlePost.PostInfo, postAttr)
			singlePost.Post.ID = int64(id)
			tpl.ExecuteTemplate(w, "postView.html", singlePost)
			singlePost.Post = nil
			singlePost.PostInfo = nil
			return
		}
		postAttr := &PostAttr{}
		user, _ := m.db.GetUserByCookie(cookie.Value)
		singlePost.AuthUser = user
		singlePost.Post, err = m.db.GetPostByPID(int64(id))
		if err != nil {
			fmt.Println("Error on PostView() function")
			WarnMessage(w, "The post not found")
			// http.Error(w, "The post not found", http.StatusNotFound)
			return
		}

		postAttr.Comments, err = m.db.GetCommentsOfPost(int64(id))
		commentRate := model.NewCommentRating()
		for i := 0; i < len(postAttr.Comments); i++ {
			commentRate = m.db.GetRateCountOfComment(postAttr.Comments[i].CommentID, int64(id))
			postAttr.Comments[i].LikeCnt = commentRate.LikeCount
			postAttr.Comments[i].DislikeCnt = commentRate.DislikeCount
		}
		postAttr.Threads, _ = m.db.GetThreadOfPost(int64(id))
		singlePost.PostInfo = append(singlePost.PostInfo, postAttr)
		singlePost.Post.ID = int64(id)
		if r.Method == "POST" {
			r.ParseForm()
			comment := model.NewComment()
			comment.Content = r.PostFormValue("comment")
			comment.CreationDate = time.Now().Format("January 2 15:04")
			comment.PostID = int64(id)
			comment.Author = user.Username
			if ok := m.db.AddComment(comment); !ok {
				fmt.Println("AddComment error")
				WarnMessage(w, "Something went wrong")
				// http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}

			postID := strconv.Itoa(id)
			http.Redirect(w, r, "/post?id="+postID, http.StatusSeeOther)
		} else {
			tpl.ExecuteTemplate(w, "postView.html", singlePost)

		}
		singlePost.Post = nil
		singlePost.PostInfo = nil
		singlePost.AuthUser = nil
	}
}

// GetAllPosts ...
func (m *Multiplexer) GetAllPosts(w http.ResponseWriter) []*model.Post {

	posts, err := m.db.GetPosts()
	if err != nil {
		WarnMessage(w, "Something went wrong")
		// http.Error(w, "Something went wrong (Test Post)", http.StatusInternalServerError)
		return nil
	}
	return posts
}
