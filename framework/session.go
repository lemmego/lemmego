package framework

import "github.com/alexedwards/scs/v2"

type Session struct {
	*scs.SessionManager
}

func NewSessionManager() *Session {
	s := scs.New()
	s.Cookie.Persist = true
	return &Session{s}
}
