package database

import (
	"net/http"

	"github.com/astgot/forum/internal/model"
)

// InsertSession ...
func (d *Database) InsertSession(u *model.Users, session *http.Cookie) (*model.Sessions, error) {
	if err := d.Open(); err != nil {
		return nil, err
	}
	cookie := model.NewSession()
	stmnt, err := d.db.Prepare("INSERT INTO Sessions (userID, cookieName, cookieValue) VALUES (?, ?, ?)")
	_, err = stmnt.Exec(u.ID, session.Name, session.Value)
	if err != nil {
		return nil, err
	}
	cookie.SessionName = session.Name
	cookie.SessionValue = session.Value
	cookie.UserID = u.ID
	return cookie, nil
}

// DeleteCookieFromDB ...
func (d *Database) DeleteCookieFromDB(cookieValue string) error {
	stmnt, err := d.db.Prepare("DELETE FROM Sessions WHERE cookieValue = ?")
	defer stmnt.Close()
	stmnt.Exec(cookieValue)
	if err != nil {
		return err
	}
	return nil

}

// GetUserByCookie ...
func (d *Database) GetUserByCookie(cookieValue string) (*model.Users, error) {
	var userID int64
	if err := d.db.QueryRow("SELECT userID from Sessions WHERE cookieValue = ?", cookieValue).Scan(&userID); err != nil {
		return nil, err
	}
	u, err := d.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return u, nil
}

/*When cookie session inserts into DB, userID repeats in the table Sessions, need to delete his old record */

// IsUserAuthenticated ...
func (d *Database) IsUserAuthenticated(u *model.Users) error {
	var cookieValue string
	if err := d.db.QueryRow("SELECT cookieValue FROM Sessions WHERE userID = ?", u.ID).Scan(&cookieValue); err != nil {
		return nil
	}
	if err := d.DeleteCookieFromDB(cookieValue); err != nil {
		return err
	}
	return nil
}
