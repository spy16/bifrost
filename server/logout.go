package server

import (
	"net/http"
)

func (srv *Server) logoutHandler(wr http.ResponseWriter, req *http.Request) {
	s := getSession(req.Context())
	if s == nil {
		writeSessionCookie(wr, "", true)
	} else {
		writeSessionCookie(wr, s.id, true)
	}
	http.Redirect(wr, req, "/login", http.StatusFound)
}
