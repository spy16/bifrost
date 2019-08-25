package server

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

// New initializes the server with given options.
func New(cfg Config, templatesDir, staticDir string) *Server {
	fsServer := newSafeFileSystemServer(staticDir)

	srv := &Server{
		cfg: cfg,
		r:   mux.NewRouter(),
		sessionStore: &sessionStore{
			sessions: map[string]*session{},
		},
		tpl: template.Must(template.ParseGlob(filepath.Join(templatesDir, "*.html"))),
	}

	// static assets
	srv.r.PathPrefix("/static").Handler(http.StripPrefix("/static", fsServer))
	srv.r.Handle("/favicon.ico", fsServer)

	// Web app handlers
	srv.r.HandleFunc("/", withSession(srv.getHomePage, srv.sessionStore))
	srv.r.HandleFunc("/login", srv.getLoginPage).Methods(http.MethodGet)
	srv.r.HandleFunc("/logout", srv.logoutHandler).Methods(http.MethodGet)
	srv.r.HandleFunc("/authorize", srv.handleAuthorize).Methods(http.MethodGet)
	srv.r.HandleFunc("/config", srv.getConfigPage).Methods(http.MethodGet)
	srv.r.HandleFunc("/callback", srv.handleCallback)

	return srv
}

// Server represents the HTTP Web app server.
type Server struct {
	cfg Config
	r   *mux.Router
	tpl *template.Template
	*sessionStore
}

func (srv *Server) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	srv.r.ServeHTTP(wr, req)
}
