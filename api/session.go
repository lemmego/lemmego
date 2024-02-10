package api

import "github.com/alexedwards/scs/v2"

type Session struct {
	*scs.SessionManager
}

func NewSessionManager() *Session {
	s := scs.New()
	s.Cookie.Persist = false
	return &Session{s}
}
