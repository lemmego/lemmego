package app

import (
	"fmt"

	"github.com/alexedwards/scs/redisstore"
	"github.com/gomodule/redigo/redis"
	"github.com/lemmego/lemmego/api/config"
	"github.com/lemmego/lemmego/api/session"
)

type SessionServiceProvider struct {
	*BaseServiceProvider
}

func (provider *SessionServiceProvider) Register(app *App) {
	// Establish connection pool to Redis.
	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%d", config.Get[string]("db.redisHost"), config.Get[int]("db.redisPort")))
		},
	}
	//sm := session.NewSession(session.NewFileStore(""))
	sm := session.NewSession(redisstore.New(pool))
	app.SetSession(sm)
}

func (provider *SessionServiceProvider) Boot() {
	//
}
