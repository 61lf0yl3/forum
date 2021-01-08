package controller

import (
	"html/template"
	"net/http"

	"github.com/astgot/forum/internal/database"
)

var tpl = template.Must(template.ParseGlob("web/templates/*"))

// Warning ...
var Warning struct {
	Warn string
}

// Multiplexer ....
type Multiplexer struct {
	Mux *http.ServeMux
	db  *database.Database
}

// NewMux ...
func NewMux() *Multiplexer {
	return &Multiplexer{
		Mux: http.NewServeMux(),
		db:  database.NewDB(database.NewConfig()),
	}
}

// WarnMessage ...
func WarnMessage(w http.ResponseWriter, warn string) {
	Warning.Warn = warn
	tpl.ExecuteTemplate(w, "error.html", Warning)
}
