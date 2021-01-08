package model

//Users ...
type Users struct {
	ID           int64  `db:"id"`
	Firstname    string `db:"firstname"`
	Lastname     string `db:"lastname"`
	Username     string `db:"username"`
	Email        string `db:"email"`
	Password     string
	ConfirmPwd   string
	EncryptedPwd string `db:"password"`
	Errors       map[string]string
}

// NewUser ...
func NewUser() *Users {
	return &Users{}
}
