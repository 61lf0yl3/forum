package controller

import (
	"regexp"
	"strings"

	"github.com/astgot/forum/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// RegExp patterns to validate user's input
var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var rxUname = regexp.MustCompile("^[\\w]")
var rxPwd = regexp.MustCompile("^[#\\w@-]{8,20}")

// ValidateInput ...
func ValidateInput(u *model.Users) bool {
	var matchL, matchF bool = true, true
	u.Errors = make(map[string]string)
	// Validity of Email
	matchEmail := rxEmail.Match([]byte(u.Email))
	if !matchEmail {
		u.Errors["Email"] = "Please enter a valid e-mail address, e.g. yourmail@example.com"
	}
	// Username
	matchUname := rxUname.Match([]byte(u.Username))
	if !matchUname {
		u.Errors["Username"] = "Please use latin alphabet, lowercase and uppercase characters with numbers or special characters"
	}
	// Password
	matchPwd := rxPwd.Match([]byte(u.Password))
	if !matchPwd {
		u.Errors["Password"] = "The password must include one lowercase, uppercase, special characters and number"
	}
	// First and Last name
	if u.Firstname != "" && u.Lastname != "" {
		matchF = rxUname.Match([]byte(u.Firstname))
		matchL = rxUname.Match([]byte(u.Lastname))
		if !matchF || !matchL {
			u.Errors["FLName"] = "Please use latin alphabet [a-z, A-Z]"
		}
	}
	if u.Password != u.ConfirmPwd {
		u.Errors["Confirm"] = "The password confirmation does not match"
	}

	return len(u.Errors) == 0
}

// HashPassword ...
func HashPassword(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	return string(hash)
}

// UnameOrEmail -->
func UnameOrEmail(query string) bool {
	if ok := strings.Contains(query, "@"); !ok {
		return false
	}
	return true
}

// ComparePassword --> Decrypt Password
func ComparePassword(hashpwd, pwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashpwd), []byte(pwd)); err != nil {
		return false
	}
	return true
}
