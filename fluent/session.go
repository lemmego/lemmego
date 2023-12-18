package fluent

import "github.com/alexedwards/scs/v2"

type Session struct {
	*scs.SessionManager
}

func NewSessionManager() *Session {
	return &Session{scs.New()}
}
