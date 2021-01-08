package database

import (
	"fmt"

	"github.com/astgot/forum/internal/model"
)

// AddRateToComment ...
func (d *Database) AddRateToComment(l *model.CommentRating, uid int64) bool {
	rate := d.GetRateCountOfComment(l.CommentID, l.PostID)
	var kind int64 // like or dislike
	if l.Like {
		kind = 1
		if rate.LikeCount == 0 && rate.DislikeCount == 0 {
			stmnt, err := d.db.Prepare("INSERT INTO CommentRating (commentID, postID, likeCount, dislikeCount) VALUES (?, ?, ?, ?)")
			defer stmnt.Close()
			_, err = stmnt.Exec(l.CommentID, l.PostID, rate.LikeCount+1, rate.DislikeCount)
			if err != nil {
				fmt.Println("db Insert CommenttRating error", err.Error())
				return false
			}
		} else {
			// Update column "likeCount" in the table
			if ok := d.AddLikeToComment(rate.LikeCount, l.CommentID, l.PostID); !ok {
				return false
			}
		}
	} else {
		kind = 0
		if rate.LikeCount == 0 && rate.DislikeCount == 0 {
			stmnt, err := d.db.Prepare("INSERT INTO CommentRating (commentID, postID, likeCount, dislikeCount) VALUES (?, ?, ?, ?)")
			defer stmnt.Close()
			_, err = stmnt.Exec(l.CommentID, l.PostID, rate.LikeCount, rate.DislikeCount+1)
			if err != nil {
				fmt.Println("db Insert CommenttRating error", err.Error())
				return false
			}
		} else {
			// Update column "dislikeCount" in the table
			if ok := d.AddDislikeToComment(rate.DislikeCount, l.CommentID, l.PostID); !ok {
				return false
			}
		}
	}

	stmnt, err := d.db.Prepare("INSERT INTO RateUserComment (commentID, postID, userID, kind) VALUES (?, ?, ?, ?)")
	defer stmnt.Close()
	_, err = stmnt.Exec(l.CommentID, l.PostID, uid, kind)
	if err != nil {
		fmt.Println("RateUserComment error")
		return false
	}
	return true
}

// UpdateRateOfComment ...
func (d *Database) UpdateRateOfComment(rate *model.CommentRating, uid int64) bool {

	// 1) What user did now? (like or dislike)
	// 2) What user have done before?
	// if user 1) liked and 2) liked ---> Delete like from post
	// If 1) liked 2) dislike ---> Delete like and add dislike
	// If 1) disliked 2)disliked ---> Delete dislike
	// If 1) disliked 2) liked ---> Delete dislike and add like

	var before int64
	if err := d.db.QueryRow("SELECT kind FROM RateUserComment WHERE userID=? AND postID=? AND commentID=?", uid, rate.PostID, rate.CommentID).
		Scan(&before); err != nil {
		fmt.Println("Select kind rate of comment error type")
		return false
	}
	rateCount := d.GetRateCountOfComment(rate.CommentID, rate.PostID)
	rate.DislikeCount = rateCount.DislikeCount
	rate.LikeCount = rateCount.LikeCount

	// Scenarios
	if before == 1 && rate.Like {
		// delete like
		if ok := d.DeleteLikeFromComment(rate.LikeCount, rate.CommentID, rate.PostID); !ok {
			return false
		}
		// delete row from RateUserPost
		if ok := d.DeleteCommentRateFromDB(uid, rate.CommentID, rate.PostID); !ok {
			return false
		}

	} else if before == 1 && !(rate.Like) {
		//delete like, add dislike
		if ok := d.DeleteLikeFromComment(rate.LikeCount, rate.CommentID, rate.PostID); !ok {
			return false
		}
		if ok := d.AddDislikeToComment(rate.DislikeCount, rate.CommentID, rate.PostID); !ok {
			return false
		}
		// Update kind equal to '0'
		if ok := d.EditCommentRate(before, uid, rate.CommentID, rate.PostID); !ok {
			return false
		}

	} else if before == 0 && !(rate.Like) {
		//delete dislike
		if ok := d.DeleteDislikeFromComment(rate.DislikeCount, rate.CommentID, rate.PostID); !ok {
			return false
		}
		// // delete row from RateUserPost
		if ok := d.DeleteCommentRateFromDB(uid, rate.CommentID, rate.PostID); !ok {
			return false
		}
	} else if before == 0 && rate.Like {
		//delete dislike, add like
		if ok := d.DeleteDislikeFromComment(rate.DislikeCount, rate.CommentID, rate.PostID); !ok {
			return false
		}
		if ok := d.AddLikeToComment(rate.LikeCount, rate.CommentID, rate.PostID); !ok {
			return false
		}
		// Update kind equal to '1'
		if ok := d.EditCommentRate(before, uid, rate.CommentID, rate.PostID); !ok {
			return false
		}
	}

	return true

}
