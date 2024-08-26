package app

import (
	"github.com/lemmego/lemmego/api/session"
)

type SessionServiceProvider struct {
	*BaseServiceProvider
}

func (provider *SessionServiceProvider) Register(app *App) {
	// Establish connection pool to Redis.
	// pool := &redis.Pool{
	// 	MaxIdle: 10,
	// 	Dial: func() (redis.Conn, error) {
	// 		conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", config.Get[string]("redis.connections.default.host"), config.Get[int]("redis.connections.default.port")))
	// 		if err != nil {
	// 			return nil, fmt.Errorf("failed to connect to redis: %v", err)
	// 		}
	// 		return conn, err
	// 	},
	// }
	// sm := session.NewSession(redisstore.New(pool))
	sm := session.NewSession(session.NewFileStore(""))
	app.SetSession(sm)
}

func (provider *SessionServiceProvider) Boot() {
	//
}
