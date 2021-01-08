package database

import (
	"fmt"

	"github.com/astgot/forum/internal/model"
)

// Create ---> signing up
func (d *Database) Create(u *model.Users) (*model.Users, error) {
	if err := d.Open(); err != nil {
		return nil, err
	}

	stmnt, err := d.db.Prepare("INSERT INTO Users (firstname, lastname, username, email, password) VALUES (?, ?, ?, ?, ?)")
	res, err := stmnt.Exec(u.Firstname, u.Lastname, u.Username, u.Email, u.EncryptedPwd)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	// assign UserID to model 'User'
	id, _ := res.LastInsertId()
	u.ID = id

	return u, nil

}

// FindByUsername ...
func (d *Database) FindByUsername(username string) (*model.Users, error) {
	if err := d.Open(); err != nil {
		return nil, err
	}

	u := &model.Users{}
	if err := d.db.QueryRow("SELECT id, username, password FROM Users where username = ?", username).Scan(
		&u.ID,
		&u.Username,
		&u.EncryptedPwd,
	); err != nil {
		return nil, err
	}

	return u, nil
}

// FindByEmail ...
func (d *Database) FindByEmail(email string) (*model.Users, error) {
	if err := d.Open(); err != nil {
		return nil, err
	}
	u := &model.Users{}
	if err := d.db.QueryRow("SELECT id, email, password FROM Users where email = ?", email).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPwd,
	); err != nil {
		fmt.Println("FindByEmail() error")
		return nil, err
	}

	return u, nil
}

// GetUserID ...
func (d *Database) GetUserID(user *model.Users, email bool) int64 {
	if err := d.Open(); err != nil {
		return 0
	}

	if email {
		if err := d.db.QueryRow("SELECT id FROM Users where email = ?", user.Username).Scan(
			&user.ID,
		); err != nil {
			return 0
		}
	} else {
		if err := d.db.QueryRow("SELECT id FROM Users where username = ?", user.Username).Scan(
			&user.ID,
		); err != nil {
			return 0
		}

	}
	return user.ID

}

// FindByUserID ...
func (d *Database) FindByUserID(UID int64) (*model.Users, error) {
	u := model.NewUser()
	if err := d.db.QueryRow("SELECT id, firstname, lastname, username FROM Users WHERE id = ?", UID).
		Scan(&u.ID,
			&u.Firstname,
			&u.Lastname,
			&u.Username); err != nil {
		return nil, err
	}
	return u, nil
}
