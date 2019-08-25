package server

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

func withSession(next http.HandlerFunc, store *sessionStore) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("sid")
		if err != nil {
			http.Redirect(wr, req, "/login", http.StatusFound)
			return
		}

		s, ok := store.sessions[cookie.Value]
		if !ok {
			writeSessionCookie(wr, "", true)
			http.Redirect(wr, req, "/login", http.StatusFound)
			return
		}
		s.id = cookie.Value

		req = req.WithContext(context.WithValue(req.Context(), sessionKey, s))
		next.ServeHTTP(wr, req)
	}
}

func getSession(ctx context.Context) *session {
	v, ok := ctx.Value(sessionKey).(*session)
	if !ok {
		return nil
	}

	return v
}

func writeSessionCookie(wr http.ResponseWriter, sid string, delete bool) {
	cookie := &http.Cookie{
		Name:     "sid",
		HttpOnly: true,
		Value:    sid,
		Expires:  time.Now().Add(24 * time.Hour),
	}

	if delete {
		cookie.Value = ""
		cookie.Expires = time.Unix(0, 0)
	}

	http.SetCookie(wr, cookie)
}

type sessionStore struct {
	sessions map[string]*session
}

type session struct {
	id    string
	Code  string
	Token *oauth2.Token
}

type contextKey string

const sessionKey = contextKey("session")
