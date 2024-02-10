package api

import (
	"fmt"
	"os"

	"github.com/alexedwards/scs/redisstore"
	"github.com/gomodule/redigo/redis"
)

type SessionServiceProvider struct {
	BaseServiceProvider
}

func (provider *SessionServiceProvider) Register(app *App) {
	// Establish connection pool to Redis.
	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")))
		},
	}
	sm := NewSessionManager()
	sm.Store = redisstore.New(pool)
	app.session = sm
}

func (provider *SessionServiceProvider) Boot() {
	//
}
