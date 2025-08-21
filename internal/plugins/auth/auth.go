package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/lemmego/api/app"
	"github.com/lemmego/api/session"
	"golang.org/x/crypto/bcrypt"
	"reflect"
)

const UserSessionKey = "user"

type Provider struct{}

type Auth struct {
	sess *session.Session
}

func New() *Auth {
	return &Auth{}
}

func (ap *Provider) Provide(a app.App) error {
	fmt.Println("Registering Auth")
	auth := &Auth{a.Session()}
	a.AddService(auth)
	return nil
}

func (a *Auth) Middleware(c app.Context) error {
	if !a.Check(c.RequestContext()) {
		return c.Unauthorized(errors.New("unauthorized"))
	}
	return c.Next()
}

func (a *Auth) Check(ctx context.Context) bool {
	if a.sess.Get(ctx, UserSessionKey) != nil {
		return true
	}

	return false
}

func (a *Auth) Login(ctx context.Context, userProvider UserProvider, username, password string) bool {
	if userProvider == nil || userProvider.GetUsername() != username {
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userProvider.GetPassword()), []byte(password)); err != nil {
		return false
	}

	a.sess.Put(ctx, UserSessionKey, userProvider)

	return true
}

func (a *Auth) Logout(ctx context.Context) {
	a.sess.Pop(ctx, UserSessionKey)
}

func Get(a app.App) *Auth {
	return a.Service(reflect.TypeOf(&Auth{})).(*Auth)
}
