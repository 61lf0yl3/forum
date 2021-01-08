package database

import (
	"fmt"

	"github.com/astgot/forum/internal/model"
)

// GetRateCountOfPost ...
func (d *Database) GetRateCountOfPost(postID int64) *model.PostRating {
	rates := model.NewPostRating()
	if err := d.db.QueryRow("SELECT * FROM PostRating WHERE postID = ?", postID).
		Scan(&rates.PostID,
			&rates.LikeCount,
			&rates.DislikeCount,
		); err != nil {
		// It means nobody rated the post, likeCount and dislikeCount now is zero
		rates.LikeCount = 0
		rates.DislikeCount = 0
	}
	return rates
}

// IsUserRatePost ...
func (d *Database) IsUserRatePost(uid, pid int64) bool {
	var comp int64
	// need to check all rated posts of the user
	res, err := d.db.Query("SELECT postID FROM RateUserPost WHERE userID=?", uid, pid)
	if err != nil {
		return false
	}
	defer res.Close()
	// Check postID is rated or not
	for res.Next() {
		if err := res.Scan(&comp); err != nil {
			fmt.Println(err.Error(), "IsUserRatePost")
			// return false
		}
		if comp == pid {
			return true
		}
		comp = 0
	}
	return false
}

// AddLike ...
func (d *Database) AddLike(likeCnt, postID int64) bool {
	stmnt, err := d.db.Prepare("UPDATE PostRating SET likeCount=? WHERE postID=?")
	defer stmnt.Close()
	_, err = stmnt.Exec(likeCnt+1, postID)
	if err != nil {
		fmt.Println("update likecount error")
		return false
	}
	return true
}

// DeleteLike ...
func (d *Database) DeleteLike(likeCnt, postID int64) bool {
	stmnt, err := d.db.Prepare("UPDATE PostRating SET likeCount=? WHERE postID=?")
	defer stmnt.Close()
	_, err = stmnt.Exec(likeCnt-1, postID)
	if err != nil {
		fmt.Println("update likecount error")
		return false
	}
	return true
}

// DeleteRateFromDB ...
func (d *Database) DeleteRateFromDB(uid, pid int64) bool {
	stmnt, err := d.db.Prepare("DELETE FROM RateUserPost WHERE userID=? AND postID=?")
	defer stmnt.Close()
	_, err = stmnt.Exec(uid, pid)
	if err != nil {
		return false
	}
	// check posts rateCounts, if [0, 0] delete row from PostRating
	rates := d.GetRateCountOfPost(pid)
	if rates.DislikeCount == 0 && rates.LikeCount == 0 {
		stmnt, err := d.db.Prepare("DELETE FROM PostRating WHERE postID=?")
		defer stmnt.Close()
		_, err = stmnt.Exec(pid)
		if err != nil {
			return false
		}
	}

	return true
}

// UpdateRate ...
func (d *Database) UpdateRate(before, uid, pid int64) bool {
	// change to like
	if before == 0 {

		stmnt, err := d.db.Prepare("UPDATE RateUserPost SET kind=1 WHERE postID=? AND userID=?")
		defer stmnt.Close()
		_, err = stmnt.Exec(pid, uid)
		if err != nil {
			fmt.Println("update likecount error")
			return false
		}
		// change to dislike
	} else {
		stmnt, err := d.db.Prepare("UPDATE RateUserPost SET kind=0 WHERE postID=? AND userID=?")
		defer stmnt.Close()
		_, err = stmnt.Exec(pid, uid)
		if err != nil {
			fmt.Println("update likecount error")
			return false
		}
	}
	return true
}

// AddDislike ...
func (d *Database) AddDislike(dislikeCnt, postID int64) bool {
	stmnt, err := d.db.Prepare("UPDATE PostRating SET dislikeCount=? WHERE postID=?")
	defer stmnt.Close()
	_, err = stmnt.Exec(dislikeCnt+1, postID)
	if err != nil {
		fmt.Println("update dislikecount error")
		return false
	}
	return true
}

// DeleteDislike ...
func (d *Database) DeleteDislike(dislikeCnt, postID int64) bool {
	stmnt, err := d.db.Prepare("UPDATE PostRating SET dislikeCount=? WHERE postID=?")
	defer stmnt.Close()
	_, err = stmnt.Exec(dislikeCnt-1, postID)
	if err != nil {
		fmt.Println("update dislikecount error")
		return false
	}
	return true
}
