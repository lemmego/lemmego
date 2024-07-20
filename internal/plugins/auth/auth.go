package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lemmego/api"
	"lemmego/api/db"
	"lemmego/api/session"
	pluginCmd "lemmego/internal/plugins/auth/cmd"
	"log"
	"net/http"
	"os"
	"strings"

	"dario.cat/mergo"

	"github.com/golang-jwt/jwt/v4"
	"github.com/invopop/validation"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

const Namespace = "fluent.auth"

var (
	ErrInvalidCreds   = errors.New("invalid credentials")
	ErrUserNotFound   = errors.New("user not found")
	ErrLoginFailed    = errors.New("login failed")
	ErrNoStrategy     = errors.New("no strategy provided: either the session manager or the token config must be provided")
	ErrNoSecret       = errors.New("no secret provided: the JWT_SECRET env variable must be provided")
	ErrNoSession      = errors.New("no session provided: the session manager must be provided")
	ErrNoUserSession  = errors.New("user session doesn't exists")
	ErrInvalidJwtSign = errors.New("invalid jwt signature")
)

type Actor interface {
	Id() string
	GetUsername() string
	GetPassword() string
}

type AuthUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// type GenericUser struct {
// 	FirstName string `json:"first_name"`
// 	LastName  string `json:"last_name"`
// 	Username  string `json:"username"`
// 	Password  string `json:"password"`
// }

type ResolveUserFunc func(c *api.Context, opts *Options) (*AuthUser, *Credentials, validation.Errors)
type CreateUserFunc func(c *api.Context, opts *Options) (bool, validation.Errors)

type CustomHandlerFunc func(opts *Options) *CustomHandlers

type CustomHandlers struct {
	ShowLogin             api.Handler
	ShowLoginEndpoint     string
	ShowRegister          api.Handler
	ShowRegisterEndpoint  string
	StoreLogin            api.Handler
	StoreLoginEndpoint    string
	StoreRegister         api.Handler
	StoreRegisterEndpoint string
}

type Options struct {
	Router            *api.Router
	DB                *db.DB
	DBFunc            func() db.DB
	Session           *session.Session
	TokenConfig       *TokenConfig
	ResolveUser       ResolveUserFunc
	CreateUser        CreateUserFunc
	GoogleOAuthConfig *oauth2.Config
	CustomViewMap     map[string]string
	HomeRoute         string
	CustomHandlers    *CustomHandlers
}

type AuthPlugin struct {
	Opts     *Options
	AuthUser *AuthUser
}

type OptFunc func(opts *Options)

type TokenConfig struct {
	Claims jwt.RegisteredClaims
}

func DefaultOptions() *Options {
	return &Options{
		HomeRoute: "/",
	}
}

func WithSessionManager(session *session.Session) OptFunc {
	return func(opts *Options) {
		opts.Session = session
	}
}

func WithTokenConfig(tokenConfig *TokenConfig) OptFunc {
	return func(opts *Options) {
		opts.TokenConfig = tokenConfig
	}
}

func New(opts ...OptFunc) *AuthPlugin {
	o := DefaultOptions()

	for _, opt := range opts {
		opt(o)
	}

	//if o.TokenConfig == nil && o.Session == nil {
	//	panic(ErrNoStrategy)
	//}

	if o.TokenConfig != nil && os.Getenv("JWT_SECRET") == "" {
		panic(ErrNoSecret)
	}

	return &AuthPlugin{o, nil}
}

func (authn *AuthPlugin) Login(ctx context.Context, a *AuthUser, username string, password string) (token string, err error) {
	// If the username and password are empty, return an error
	if a.Username == "" && a.Password == "" {
		return "", ErrInvalidCreds
	}

	// If the username doesn't match the one provided, return an error
	if a.Username != username {
		return "", ErrUserNotFound
	}

	// If the password doesn't match the one provided, return an error
	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)); a.Username != "" && a.Password != "" && err != nil {
		return "", ErrInvalidCreds
	}

	// If the session manager is provided, store the user in the session
	if authn.Opts.Session != nil {
		userJson, _ := json.Marshal(a)
		authn.Opts.Session.Put(ctx, "user", string(userJson))
		authn.Opts.Session.Put(ctx, "userId", a.ID)
	} else {
		return "", ErrNoSession
	}

	// If the token config is provided, generate a token
	if authn.Opts.TokenConfig != nil {
		mergo.Merge(&authn.Opts.TokenConfig.Claims, jwt.RegisteredClaims{
			Subject: a.ID,
		})
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, authn.Opts.TokenConfig.Claims)

		token, err = claims.SignedString([]byte(os.Getenv("JWT_SECRET")))

		if err != nil {
			return "", ErrInvalidJwtSign
		}
	}

	// Return the token and error
	return token, err
}

func (authn *AuthPlugin) ForceLogin(ctx context.Context, a Actor) (token string, err error) {
	if a.GetUsername() != "" && a.GetPassword() != "" {
		if authn.Opts.Session != nil {
			authn.Opts.Session.Put(ctx, "userId", a.Id())
		}
		if authn.Opts.TokenConfig != nil {
			mergo.Merge(&authn.Opts.TokenConfig.Claims, jwt.RegisteredClaims{
				Subject: a.Id(),
			})
			claims := jwt.NewWithClaims(jwt.SigningMethodHS256, authn.Opts.TokenConfig.Claims)

			token, err = claims.SignedString([]byte(os.Getenv("JWT_SECRET")))
		}
	} else {
		return "", ErrLoginFailed
	}

	return token, err
}

func (authn *AuthPlugin) Check(r *http.Request) error {
	user := &AuthUser{}
	if authn.Opts.Session != nil {
		if exists := authn.Opts.Session.Exists(r.Context(), "userId"); exists {
			json.Unmarshal([]byte(authn.Opts.Session.Get(r.Context(), "user").(string)), user)
			authn.AuthUser = user
			return nil
		} else {
			return ErrNoUserSession
		}
	}

	if authn.Opts.TokenConfig != nil {
		jwtToken := ""
		jwtCookie, err := r.Cookie("jwt")
		if err == nil {
			jwtToken = strings.Replace(jwtCookie.Value, "jwt=", "", -1)
		} else {
			jwtToken = strings.Replace(r.Header.Get("Authorization"), "bearer ", "", -1)
		}
		if jwtToken == "" {
			return errors.New("jwt cookie not found")
		}
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			return err
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			user.ID = claims["sub"].(string)
			user.Username = claims["username"].(string)
		}

		authn.AuthUser = user
	}
	return errors.New("could not parse jwt")
}

// Guard the route with the auth middleware
func (authn *AuthPlugin) Guard(next api.Handler) api.Handler {
	return func(c *api.Context) error {
		if err := authn.Check(c.Request()); err != nil {
			return c.Respond(http.StatusUnauthorized, &api.R{
				Payload:    api.M{"message": "Unauthorized"},
				RedirectTo: "/login",
			})
		} else {
			c.Set("user", authn.AuthUser)
			return next(c)
		}
	}
}

// Disallow authenticated users from accessing a route
func (authn *AuthPlugin) Guest(next api.Handler) api.Handler {
	return func(c *api.Context) error {
		if err := authn.Check(c.Request()); err == nil {
			return c.Respond(http.StatusUnauthorized, &api.R{
				Payload:    api.M{"message": "Unauthorized"},
				RedirectTo: "/",
			})
		} else {
			return next(c)
		}
	}
}

func (p *AuthPlugin) Namespace() string {
	return Namespace
}

func (p *AuthPlugin) Commands() []*cobra.Command {
	return []*cobra.Command{}
}

func (p *AuthPlugin) InstallCommand() *cobra.Command {
	return pluginCmd.GetInstallCommand(p)
}

func (p *AuthPlugin) Boot(app *api.App) error {
	p.Opts.Session = app.Session()
	p.Opts.DB = app.Db()
	p.Opts.Router = app.Router()
	return nil
}

func (p *AuthPlugin) EventListeners() map[string]func() {
	return nil
}

func (p *AuthPlugin) Migrations() []string {
	return nil
}

func (p *AuthPlugin) Templates() map[string][]byte {
	return nil
	return map[string][]byte{
		// "login.page.tmpl":    loginTmpl,
		// "register.page.tmpl": registerTmpl,
	}
}

func (p *AuthPlugin) Middlewares() []api.Middleware {
	return nil
}

func (p *AuthPlugin) RouteMiddlewares() map[string]api.Middleware {
	return map[string]api.Middleware{
		"auth": p.Guard,
	}
}

func (p *AuthPlugin) indexLoginPageHandler() api.Handler {
	return func(c *api.Context) error {
		return c.Render(200, "login.page.tmpl", nil)
	}
}

func (p *AuthPlugin) indexRegisterPageHandler() api.Handler {
	return func(c *api.Context) error {
		return c.Render(200, "register.page.tmpl", nil)
	}
}

func (p *AuthPlugin) storeRegisterHandler() api.Handler {
	return func(c *api.Context) error {
		ok, errs := p.Opts.CreateUser(c, p.Opts)

		if errs != nil {
			return c.WithErrors(errs).Back()
		}

		if ok {
			return c.Respond(http.StatusOK, &api.R{
				Message:    c.Alert("success", "Registration Successful"),
				RedirectTo: "/login",
				Payload:    api.M{"message": "Registration Successful"},
			})
		} else {
			return c.WithError("Registration Failed").Back()
		}
	}
}

func (p *AuthPlugin) storeLoginHandler() api.Handler {
	return func(c *api.Context) error {
		if aUser, creds, err := p.Opts.ResolveUser(c, p.Opts); err != nil {
			return c.Respond(http.StatusBadRequest, &api.R{
				Message:    c.Alert("error", "Login failed."),
				RedirectTo: "/login",
				Payload:    api.M{"message": "Login failed."},
			})
		} else {
			token, err := p.Login(c.Request().Context(), aUser, creds.Username, creds.Password)

			if err != nil {
				log.Println(err)
				payload := api.M{"message": "Login failed."}
				if errors.Is(err, ErrInvalidCreds) {
					payload["message"] = "Invalid credentials"
				}

				if errors.Is(err, ErrUserNotFound) {
					payload["message"] = "User not found"
				}

				return c.Respond(http.StatusBadRequest, &api.R{
					Message:    c.Alert("error", payload["message"].(string)),
					RedirectTo: "/login",
					Payload:    payload,
				})
			}

			payload := api.M{"message": "Login successful."}

			if token != "" {
				c.SetCookie("jwt", token, 60*60*24*7, "/", "", false, true)
				payload["token"] = token
			}

			return c.Respond(http.StatusOK, &api.R{
				Message:    c.Alert("success", "Login successful."),
				Payload:    payload,
				RedirectTo: p.Opts.HomeRoute,
			})
		}
	}
}

func (p *AuthPlugin) Routes() []*api.Route {
	routes := []*api.Route{
		&api.Route{
			Method: http.MethodGet,
			Path:   "/login",
			Handler: func(c *api.Context) error {
				return c.Inertia("Forms/Login", nil)
			},
		},
		&api.Route{
			Method: http.MethodGet,
			Path:   "/register",
			Handler: func(c *api.Context) error {
				return c.Inertia("Forms/Register", nil)
			},
		},
	}

	return routes
}

func (p *AuthPlugin) Webhooks() []string {
	return nil
}
