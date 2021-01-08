package controller

import (
	"fmt"
	"net/http"

	"github.com/astgot/forum/internal/model"
)

// MainHandle ...
func (m *Multiplexer) MainHandle() http.HandlerFunc {

	// Need to create structure to show array of Users, Posts, Comments, Categories for arranging them in HTML
	type PostRaw struct {
		Post     *model.Post
		Threads  []*model.Thread
		PostRate *model.PostRating
		// Comments []*model.Comments
	}
	var mainPage struct {
		AuthUser   *model.Users
		PostScroll []*PostRaw
	}
	// Here we can create our own struct, which is usable only here
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/main" {
			WarnMessage(w, "404 Not Found")
			// http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}
		posts := m.GetAllPosts(w)

		cookie, err := r.Cookie("authenticated")
		if err != nil {
			// if user is guest, retrieve all posts for displaying
			for _, post := range posts {
				guest := &PostRaw{}
				// guest.Comments, _ = m.db.GetCommentsOfPost(post.PostID)
				// guest.Author, _ = m.db.FindByUserID(post.UserID)
				guest.Post = post
				guest.Threads, _ = m.db.GetThreadOfPost(post.ID)
				guest.PostRate = m.db.GetRateCountOfPost(post.ID)
				mainPage.PostScroll = append(mainPage.PostScroll, guest)
			}
			tpl.ExecuteTemplate(w, "main.html", mainPage)
			mainPage.PostScroll = nil
			return
		}
		// if User is authenticated
		user, _ := m.db.GetUserByCookie(cookie.Value)
		mainPage.AuthUser = user
		for _, post := range posts {
			auth := &PostRaw{}
			// auth.Comments, _ = m.db.GetCommentsOfPost(post.PostID)
			auth.Post = post
			auth.Threads, err = m.db.GetThreadOfPost(post.ID)
			if err != nil {
				WarnMessage(w, "Something went wrong")
				fmt.Println("Threads retrieving error")
				return
			}
			auth.PostRate = m.db.GetRateCountOfPost(post.ID)
			mainPage.PostScroll = append(mainPage.PostScroll, auth)
		}
		tpl.ExecuteTemplate(w, "main.html", mainPage)
		// prevent from posts doubling
		mainPage.AuthUser = nil
		mainPage.PostScroll = nil
	}
}
