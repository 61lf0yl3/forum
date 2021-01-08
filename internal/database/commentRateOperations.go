package database

import (
	"fmt"

	"github.com/astgot/forum/internal/model"
)

// GetRateCountOfComment ...
func (d *Database) GetRateCountOfComment(cid, pid int64) *model.CommentRating {
	rate := model.NewCommentRating()
	if err := d.db.QueryRow("SELECT likeCount, dislikeCount FROM CommentRating WHERE commentID=? AND postID=?", cid, pid).
		Scan(&rate.LikeCount,
			&rate.DislikeCount,
		); err != nil {
		// It means nobody rated the comment, likeCount and dislikeCount now is zero
		rate.LikeCount = 0
		rate.DislikeCount = 0
	}
	return rate
}

// IsUserRateComment ...
func (d *Database) IsUserRateComment(uid, pid, cid int64) bool {
	var comp int64
	query, err := d.db.Query("SELECT commentID FROM RateUserComment WHERE userID=? AND postID=?", uid, pid)
	if err != nil {
		return false
	}
	defer query.Close()
	for query.Next() {
		if err := query.Scan(&comp); err != nil {
			fmt.Println(err.Error(), "IsUserRateComment")
		}
		if cid == comp {
			return true
		}
		comp = 0
	}
	return false
}

// AddLikeToComment ..
func (d *Database) AddLikeToComment(likeCnt, cid, pid int64) bool {
	stmnt, err := d.db.Prepare("UPDATE CommentRating SET likeCount=? WHERE commentID=? AND postID=? ")
	defer stmnt.Close()
	_, err = stmnt.Exec(likeCnt+1, cid, pid)
	if err != nil {
		fmt.Println("update comment likecount error", err.Error())
		return false
	}
	return true
}

// DeleteLikeFromComment ...
func (d *Database) DeleteLikeFromComment(likeCnt, cid, pid int64) bool {
	stmnt, err := d.db.Prepare("UPDATE CommentRating SET likeCount=? WHERE commentID=? AND postID=?")
	defer stmnt.Close()
	_, err = stmnt.Exec(likeCnt-1, cid, pid)
	if err != nil {
		fmt.Println("update comment likecount error", err.Error())
		return false
	}
	return true
}

// AddDislikeToComment ...
func (d *Database) AddDislikeToComment(dislikeCnt, cid, pid int64) bool {
	stmnt, err := d.db.Prepare("UPDATE CommentRating SET dislikeCount=? WHERE commentID=? AND postID=?")
	defer stmnt.Close()
	_, err = stmnt.Exec(dislikeCnt+1, cid, pid)
	if err != nil {
		fmt.Println("update comment dislikecount error", err.Error())
		return false
	}
	return true
}

// DeleteDislikeFromComment ...
func (d *Database) DeleteDislikeFromComment(dislikeCnt, cid, pid int64) bool {
	stmnt, err := d.db.Prepare("UPDATE CommentRating SET dislikeCount=? WHERE commentID=? AND postID=?")
	defer stmnt.Close()
	_, err = stmnt.Exec(dislikeCnt-1, cid, pid)
	if err != nil {
		fmt.Println("update comment dislikecount error", err.Error())
		return false
	}
	return true
}

// DeleteCommentRateFromDB ...
func (d *Database) DeleteCommentRateFromDB(uid, cid, pid int64) bool {
	stmnt, err := d.db.Prepare("DELETE FROM RateUserComment WHERE commentID =? AND userID=? AND postID=?")
	defer stmnt.Close()
	_, err = stmnt.Exec(cid, uid, pid)
	if err != nil {
		return false
	}
	// check posts rateCounts, if [0, 0] delete row from CommentRating
	rates := d.GetRateCountOfComment(cid, pid)
	if rates.DislikeCount == 0 && rates.LikeCount == 0 {
		stmnt, err := d.db.Prepare("DELETE FROM CommentRating WHERE commentID=? AND postID=?")
		defer stmnt.Close()
		_, err = stmnt.Exec(cid, pid)
		if err != nil {
			return false
		}
	}
	return true
}

// EditCommentRate ...
func (d *Database) EditCommentRate(before, uid, cid, pid int64) bool {
	// change to like
	if before == 0 {

		stmnt, err := d.db.Prepare("UPDATE RateUserComment SET kind=1 WHERE commentID=? AND postID=? AND userID=?")
		defer stmnt.Close()
		_, err = stmnt.Exec(cid, pid, uid)
		if err != nil {
			fmt.Println("update comment likecount error", err.Error())
			return false
		}
		// change to dislike
	} else {
		stmnt, err := d.db.Prepare("UPDATE RateUserComment SET kind=0 WHERE commentID=? AND postID=? AND userID=?")
		defer stmnt.Close()
		_, err = stmnt.Exec(cid, pid, uid)
		if err != nil {
			fmt.Println("update comment likecount error", err.Error())
			return false
		}
	}
	return true
}
