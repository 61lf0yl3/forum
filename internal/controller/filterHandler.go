package controller

import (
	"net/http"

	"github.com/astgot/forum/internal/model"
)

// FilterHandler ...
/* To demonstrate:
   1) created posts --> /filter?section=my_posts
   2) liked posts   --> /filter?section=rated
   3) by categories --> /filter?search=<request>
*/
func (m *Multiplexer) FilterHandler() http.HandlerFunc {
	type PostRaw struct {
		Post     *model.Post
		Threads  []*model.Thread
		PostRate *model.PostRating
	}
	var filter struct {
		Section    string
		AuthUser   *model.Users
		PostScroll []*PostRaw
		Error      string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/filter" {
			WarnMessage(w, "404 Not Found")
			// http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}

		section := r.URL.Query().Get("section")
		if section == "my_posts" {

			c, err := r.Cookie("authenticated")
			if err != nil {
				WarnMessage(w, "You need to be authorized")
				// http.Error(w, "You need to be authenticated", http.StatusUnauthorized)
				return
			}
			user, err := m.db.GetUserByCookie(c.Value)
			if err != nil {
				WarnMessage(w, "Something went wrong")
				// http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
			// Add function to find posts of users by their ID
			/* 1) retrive posts from table "Posts"
			 */
			posts, err := m.db.GetPostsByUID(user.ID)
			if err != nil {
				WarnMessage(w, "Something went wrong")
				// http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
			for _, post := range posts {
				result := &PostRaw{}
				result.Post = post
				result.Threads, _ = m.db.GetThreadOfPost(post.ID)
				result.PostRate = m.db.GetRateCountOfPost(post.ID)
				filter.PostScroll = append(filter.PostScroll, result)
			}
			filter.Section = "My Posts"
			filter.AuthUser = user
			tpl.ExecuteTemplate(w, "filter.html", filter)
			//Prevent from doubling of content
			filter.Section = ""
			filter.AuthUser = nil
			filter.PostScroll = nil
			return
		} else if section == "liked" {
			c, err := r.Cookie("authenticated")
			if err != nil {
				WarnMessage(w, "You need to be authorized")
				// http.Error(w, "You need to be authenticated", http.StatusUnauthorized)
				return
			}
			user, err := m.db.GetUserByCookie(c.Value)
			if err != nil {
				WarnMessage(w, "Something went wrong")
				// http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
			// Add function to find post rated by the user
			/* 1) retrieve liked posts from "RateUserPost"
			 */
			posts, err := m.db.GetRatedPostsByUID(user.ID)
			if err != nil {
				WarnMessage(w, "Something went wrong")
				// http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
			for _, post := range posts {
				result := &PostRaw{}
				result.Post = post
				result.Threads, _ = m.db.GetThreadOfPost(post.ID)
				result.PostRate = m.db.GetRateCountOfPost(post.ID)
				filter.PostScroll = append(filter.PostScroll, result)

			}
			filter.Section = "Liked posts"
			filter.AuthUser = user
			tpl.ExecuteTemplate(w, "filter.html", filter)
			filter.Section = ""
			filter.AuthUser = nil
			filter.PostScroll = nil
			return
		}

		// Search by category
		r.ParseForm()
		search := r.PostFormValue("search")
		// search := r.URL.Query().Get("search")
		c, err := r.Cookie("authenticated")
		if err == nil {
			user, err := m.db.GetUserByCookie(c.Value)
			if err != nil {
				WarnMessage(w, "Something went wrong")
				// http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
			filter.AuthUser = user
		}
		filter.Section = "Search results for \"" + search + "\""
		/* 1)search that category in the table "Threads" (retrieve all threadID)
		   2)search posts in the table "PostMapping" (retrieve all postID)
		   3)show posts from table "Posts"
		*/
		posts, err := m.db.SearchThread(search)
		if err != nil {
			WarnMessage(w, "Something went wrong")
			// http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		// If search doesn't give any results
		if len(posts) == 0 {
			filter.Error = "No results for \"" + search + "\""
		} else {
			filter.Error = filter.Section
		}
		for _, post := range posts {
			result := &PostRaw{}
			result.Post = post
			result.Threads, _ = m.db.GetThreadOfPost(post.ID)
			result.PostRate = m.db.GetRateCountOfPost(post.ID)
			filter.PostScroll = append(filter.PostScroll, result)
		}
		tpl.ExecuteTemplate(w, "filter.html", filter)
		filter.AuthUser = nil
		filter.Section = ""
		filter.PostScroll = nil
		filter.Error = ""
	}
}
