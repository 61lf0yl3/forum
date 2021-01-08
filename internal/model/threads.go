package model

// Thread ...
type Thread struct {
	ID   int64  `db:"ID"`
	Name string `db:"Name"`
}

// NewThread ....
func NewThread() *Thread {
	return &Thread{}
}
