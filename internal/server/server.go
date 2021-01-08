package server

import (
	"fmt"
	"net/http"

	"github.com/astgot/forum/internal/controller"
	"github.com/astgot/forum/internal/database"
)

// Server ..
type Server struct {
	config   *Config
	database *database.Database
	mux      *controller.Multiplexer
}

// New - generates instance to support service
func New(config *Config) *Server {
	return &Server{
		config: config,
		mux:    controller.NewMux(),
	}
}

// Start - Initializing server
func (s *Server) Start() error {

	if err := s.ConfigureDB(); err != nil {
		return err
	}
	s.ConfigureRouter()

	fmt.Println("Server is working on port :8080 ...")

	return http.ListenAndServe(s.config.WebPort, s.mux.Mux)

}

// ConfigureRouter ...
func (s *Server) ConfigureRouter() {
	fs := http.FileServer(http.Dir("web/css"))
	s.mux.Mux.Handle("/css/", http.StripPrefix("/css/", fs))
	s.mux.Mux.HandleFunc("/", s.mux.MainHandle())
	s.mux.Mux.HandleFunc("/signup", s.mux.SignupHandle())
	s.mux.Mux.HandleFunc("/login", s.mux.LoginHandle())
	s.mux.Mux.HandleFunc("/logout", s.mux.LogoutHandle())
	s.mux.Mux.HandleFunc("/confirmation", controller.ConfirmHandler)
	s.mux.Mux.HandleFunc("/create", s.mux.CreatePostHandler())
	s.mux.Mux.HandleFunc("/post", s.mux.PostView())
	s.mux.Mux.HandleFunc("/rate", s.mux.RateHandler())
	s.mux.Mux.HandleFunc("/filter", s.mux.FilterHandler())
	return
}

// ConfigureDB ...
func (s *Server) ConfigureDB() error {
	db := database.NewDB(s.config.Database)
	if err := db.InitDB(); err != nil {
		return err
	}
	s.database = db //fill Server with DB instance
	return nil
}
