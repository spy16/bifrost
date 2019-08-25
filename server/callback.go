package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func (srv *Server) handleCallback(wr http.ResponseWriter, req *http.Request) {
	state := req.URL.Query().Get("state")

	stateCookie, err := req.Cookie("_s")
	if err != nil {
		renderErr(wr, err)
		return
	}

	providerCookie, err := req.Cookie("_provider")
	if err != nil {
		renderErr(wr, err)
		return
	}

	if state != stateCookie.Value {
		renderErr(wr, errors.New("state parameters mismatch"))
		return
	}

	code := strings.TrimSpace(req.URL.Query().Get("code"))
	if code == "" {
		renderErr(wr, errors.New("code cannot be empty"))
		return
	}

	provider, found := srv.cfg.Providers[providerCookie.Value]
	if !found {
		renderErr(wr, fmt.Errorf("provider '%s' not found", providerCookie.Value))
		return
	}

	token, err := provider.OAuth2().Exchange(req.Context(), code)
	if err != nil {
		renderErr(wr, err)
		return
	}

	sid := randStringRunes(10)
	srv.sessions[sid] = &session{
		Code:  code,
		Token: token,
	}
	writeSessionCookie(wr, sid, false)

	http.Redirect(wr, req, "/", http.StatusFound)
}

func renderErr(wr http.ResponseWriter, e error) {
	wr.Write([]byte(e.Error()))
}
