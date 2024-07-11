package session

import (
	"github.com/alexedwards/scs/v2"
)

type Session struct {
	*scs.SessionManager
}

func NewSession(store scs.Store) *Session {
	s := scs.New()
	s.Store = store
	s.Cookie.Persist = false
	return &Session{s}
}
