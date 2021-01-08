package database

import "log"

// BuildSchema ---> create Tables
func (d *Database) BuildSchema() error {
	users, err := d.db.Prepare(`CREATE TABLE IF NOT EXISTS Users (
									id INTEGER PRIMARY KEY NOT NULL, 
									firstname TEXT NOT NULL, 
									lastname TEXT NOT NULL, 
									username TEXT NOT NULL UNIQUE, 
									email TEXT NOT NULL UNIQUE, 
									password TEXT NOT NULL
								)`)

	defer users.Close()
	CheckErr(err)
	users.Exec()

	sessions, err := d.db.Prepare(`CREATE TABLE IF NOT EXISTS Sessions (
		userID INTEGER NOT NULL,
		cookieName TEXT NOT NULL,
		cookieValue TEXT NOT NULL UNIQUE,
		FOREIGN KEY(userID) REFERENCES Users(id)
	)`)
	defer sessions.Close()
	CheckErr(err)
	sessions.Exec()

	posts, err := d.db.Prepare(`CREATE TABLE IF NOT EXISTS Posts (
		post_id INTEGER PRIMARY KEY NOT NULL,
		user_id INTEGER NOT NULL,
		author TEXT NOT NULL, 
		title TEXT NOT NULL, 
		content TEXT NOT NULL,
		creationDate TEXT NOT NULL,
		FOREIGN KEY(user_id) REFERENCES Users(id) 
	)`)
	defer posts.Close()
	CheckErr(err)
	posts.Exec()

	threads, err := d.db.Prepare(`CREATE TABLE IF NOT EXISTS Threads (
		ID INTEGER PRIMARY KEY NOT NULL, 
		Name TEXT NOT NULL UNIQUE
	)`)
	defer threads.Close()
	CheckErr(err)
	threads.Exec()

	postMap, err := d.db.Prepare(`CREATE TABLE IF NOT EXISTS PostMapping (
		postID INTEGER NOT NULL, 
		threadID INTEGER NOT NULL,
		FOREIGN KEY(postID) REFERENCES Posts(post_id),
		FOREIGN KEY(threadID) REFERENCES Threads(ID)
	)`)
	defer postMap.Close()
	CheckErr(err)
	postMap.Exec()

	comments, err := d.db.Prepare(`CREATE TABLE IF NOT EXISTS Comments (
		ID INTEGER PRIMARY KEY NOT NULL,
		postID INTEGER NOT NULL,
		author TEXT NOT NULL,
		content TEXT NOT NULL, 
		creationDate TEXT NOT NULL,
		FOREIGN KEY(postID) REFERENCES Posts(post_id)
		)`)
	// FOREIGN KEY(userID) REFERENCES Users(id)
	defer comments.Close()
	CheckErr(err)
	comments.Exec()

	postRating, err := d.db.Prepare(`CREATE TABLE IF NOT EXISTS PostRating (
		postID INTEGER NOT NULL,
		likeCount INTEGER NOT NULL,
		dislikeCount INTEGER NOT NULL,
		FOREIGN KEY(postID) REFERENCES Posts(post_id)
	)`)
	defer postRating.Close()
	CheckErr(err)
	postRating.Exec()

	commRating, err := d.db.Prepare(`CREATE TABLE IF NOT EXISTS CommentRating (
		commentID INTEGER NOT NULL,
		postID INTEGER NOT NULL,
		likeCount INTEGER NOT NULL,
		dislikeCount INTEGER NOT NULL,
		FOREIGN KEY(commentID) REFERENCES Comments(ID),
		FOREIGN KEY(postID) REFERENCES Posts(post_id)
	)`)
	defer commRating.Close()
	CheckErr(err)
	commRating.Exec()

	rateUP, err := d.db.Prepare(`CREATE TABLE IF NOT EXISTS RateUserPost (
		userID INTEGER NOT NULL,
		postID INTEGER NOT NULL,
		kind INTEGER NOT NULL,
		FOREIGN KEY(userID) REFERENCES Users(ID),
		FOREIGN KEY(postID) REFERENCES Posts(post_id)
	)`)
	defer rateUP.Close()
	CheckErr(err)
	rateUP.Exec()

	rateUC, err := d.db.Prepare(`CREATE TABLE IF NOT EXISTS RateUserComment (
		commentID INTEGER NOT NULL,
		postID INTEGER NOT NULL,
		userID INTEGER NOT NULL,
		kind INTEGER NOT NULL,
		FOREIGN KEY(commentID) REFERENCES Comments(ID),
		FOREIGN KEY(postID) REFERENCES Posts(post_id),
		FOREIGN KEY(userID) REFERENCES Users(ID)
	)`)
	defer rateUC.Close()
	CheckErr(err)
	rateUC.Exec()
	return nil
}

// CheckErr ...
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// like_count INTEGER,
// dislike_count INTEGER,
