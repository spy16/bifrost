package server

import (
	"fmt"
	"net/http"
	"time"
)

func (srv *Server) getLoginPage(wr http.ResponseWriter, req *http.Request) {
	var loginData []providerInfo
	for name, p := range srv.cfg.Providers {
		loginData = append(loginData, providerInfo{
			Name:     name,
			ClientID: p.ClientID,
		})
	}

	srv.tpl.ExecuteTemplate(wr, "login.html", loginData)
}

func (srv *Server) handleAuthorize(wr http.ResponseWriter, req *http.Request) {
	providerName := req.FormValue("provider")

	provider, found := srv.cfg.Providers[providerName]
	if !found {
		renderErr(wr, fmt.Errorf("provider '%s' not found", providerName))
		return
	}

	state := randStringRunes(8)
	http.SetCookie(wr, &http.Cookie{
		Name:     "_s",
		Value:    state,
		Expires:  time.Now().Add(5 * time.Minute),
		HttpOnly: true,
	})
	http.SetCookie(wr, &http.Cookie{
		Name:     "_provider",
		Value:    providerName,
		Expires:  time.Now().Add(5 * time.Minute),
		HttpOnly: true,
	})
	http.Redirect(wr, req, provider.AuthorizeURL(state), http.StatusFound)
}

type providerInfo struct {
	Name     string
	ClientID string
}
