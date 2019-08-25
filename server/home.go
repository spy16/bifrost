package server

import (
	"net/http"
)

func (srv *Server) getHomePage(wr http.ResponseWriter, req *http.Request) {
	session := getSession(req.Context())
	if session == nil {
		writeSessionCookie(wr, "", true)
		http.Redirect(wr, req, "/login", http.StatusFound)
		return
	}

	srv.tpl.ExecuteTemplate(wr, "home.html", *session)
}

func (srv *Server) getConfigPage(wr http.ResponseWriter, req *http.Request) {
	srv.tpl.ExecuteTemplate(wr, "config.html", map[string]interface{}{})
}
