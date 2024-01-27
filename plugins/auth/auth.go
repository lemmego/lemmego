package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"pressebo/framework"
	"strings"

	"dario.cat/mergo"

	"github.com/golang-jwt/jwt/v4"
	"github.com/invopop/validation"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

const Namespace = "fluent.auth"

var (
	ErrLoginFailed   = errors.New("login failed")
	ErrNoStrategy    = errors.New("no strategy provided: either the session manager or the token config must be provided")
	ErrNoSecret      = errors.New("no secret provided: the JWT_SECRET env variable must be provided")
	ErrNoSession     = errors.New("no session provided: the session manager must be provided")
	ErrNoUserSession = errors.New("user session doesn't exists")
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

type ResolveUserFunc func(username string) (*AuthUser, error)
type CreateUserFunc func(firstName string, lastName string, username string, password string) bool

type Options struct {
	DB                framework.DBSession
	DBFunc            func() framework.DBSession
	Session           *framework.Session
	TokenConfig       *TokenConfig
	ResolveUser       ResolveUserFunc
	CreateUser        CreateUserFunc
	GoogleOAuthConfig *oauth2.Config
	CustomViewMap     map[string]string
	HomeRoute         string
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

func WithDefaultUserResolver(opts *Options) {
	opts.ResolveUser = func(username string) (*AuthUser, error) {
		db := opts.DB
		authUser := AuthUser{}
		q, err := db.SQL().QueryRow("select id, email, password from users where email = $1 limit 1", username)
		if err != nil {
			return nil, err
		}

		if err := q.Scan(&authUser.ID, &authUser.Username, &authUser.Password); err != nil {
			return nil, err
		}
		return &authUser, nil
	}
}

func WithDefaultUserCreator(opts *Options) {
	opts.CreateUser = func(firstName string, lastName string, username string, password string) bool {
		db := opts.DB

		encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		q := db.
			SQL().
			InsertInto("users").
			Columns("first_name", "last_name", "email", "password").
			Values(firstName, lastName, username, encryptedPassword)

		if _, err := q.Exec(); err != nil {
			return false
		}
		return true
	}
}

func WithUserResolver(resolveUser ResolveUserFunc) OptFunc {
	return func(opts *Options) {
		opts.ResolveUser = resolveUser
	}
}

func WithUserCreator(createUser CreateUserFunc) OptFunc {
	return func(opts *Options) {
		opts.CreateUser = createUser
	}
}

func WithSessionManager(session *framework.Session) OptFunc {
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
	WithDefaultUserResolver(o)
	WithDefaultUserCreator(o)

	for _, opt := range opts {
		opt(o)
	}

	// if o.TokenConfig == nil && o.Session == nil {
	// 	panic(ErrNoStrategy)
	// }

	if o.TokenConfig != nil && os.Getenv("JWT_SECRET") == "" {
		panic(ErrNoSecret)
	}

	return &AuthPlugin{o, nil}
}

func (authn *AuthPlugin) Login(ctx context.Context, a *AuthUser, username string, password string) (token string, err error) {
	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)); a.Username != "" && a.Password != "" && err == nil {
		if authn.Opts.Session != nil {
			userJson, _ := json.Marshal(a)
			authn.Opts.Session.Put(ctx, "user", string(userJson))
			authn.Opts.Session.Put(ctx, "userId", a.ID)
		} else {
			return "", ErrNoSession
		}
		if authn.Opts.TokenConfig != nil {
			mergo.Merge(&authn.Opts.TokenConfig.Claims, jwt.RegisteredClaims{
				Subject: a.ID,
			})
			claims := jwt.NewWithClaims(jwt.SigningMethodHS256, authn.Opts.TokenConfig.Claims)

			token, err = claims.SignedString([]byte(os.Getenv("JWT_SECRET")))
		}
	} else {
		return "", ErrLoginFailed
	}

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
func (authn *AuthPlugin) Guard(next framework.Handler) framework.Handler {
	return func(c *framework.Context) error {
		if err := authn.Check(c.Request()); err != nil {
			return c.Respond(&framework.R{
				Status:     http.StatusUnauthorized,
				Payload:    framework.M{"message": "Unauthorized"},
				RedirectTo: "/login",
			})
		} else {
			c.Set("user", authn.AuthUser)
			return next(c)
		}
	}
}

// Disallow authenticated users from accessing a route
func (authn *AuthPlugin) Guest(next framework.Handler) framework.Handler {
	return func(c *framework.Context) error {
		if err := authn.Check(c.Request()); err == nil {
			return c.Respond(&framework.R{
				Status:     http.StatusUnauthorized,
				Payload:    framework.M{"message": "Unauthorized"},
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

func (p *AuthPlugin) Boot(app *framework.App) error {
	p.Opts.Session = app.Session()
	p.Opts.DB = app.Db()
	return nil
}

func (p *AuthPlugin) EventListeners() map[string]func() {
	return nil
}

func (p *AuthPlugin) Migrations() []string {
	return nil
}

func (p *AuthPlugin) Templates() map[string][]byte {
	return map[string][]byte{
		"login.page.tmpl":    loginTmpl,
		"register.page.tmpl": registerTmpl,
	}
}

func (p *AuthPlugin) Middlewares() []func(http.Handler) http.Handler {
	return nil
}

func (p *AuthPlugin) RouteMiddlewares() map[string]framework.Middleware {
	return map[string]framework.Middleware{
		"auth": p.Guard,
	}
}

func (p *AuthPlugin) indexLoginPageHandler() framework.Handler {
	return func(c *framework.Context) error {
		return c.Render(200, "login.page.tmpl", nil)
	}
}

func (p *AuthPlugin) indexRegisterPageHandler() framework.Handler {
	return func(c *framework.Context) error {
		return c.Render(200, "register.page.tmpl", nil)
	}
}
func (p *AuthPlugin) storeRegisterHandler() framework.Handler {
	return func(c *framework.Context) error {
		body := c.GetBody()
		registrationRequest := &RegistrationStoreRequest{
			FirstName:            body["first_name"][0],
			LastName:             body["last_name"][0],
			Username:             body["username"][0],
			Password:             body["password"][0],
			PasswordConfirmation: body["password_confirmation"][0],
		}

		if err := c.Validate(registrationRequest); err != nil {
			return c.WithErrors(err.(validation.Errors)).Back()
		}

		ok := p.Opts.CreateUser(
			registrationRequest.FirstName,
			registrationRequest.LastName,
			registrationRequest.Username,
			registrationRequest.Password,
		)

		if ok {
			return c.Respond(&framework.R{
				Message:    &framework.AlertMessage{"success", "Registration successful."},
				RedirectTo: "/login",
			})
		} else {
			return c.WithError("Registration Failed").Back()
		}
	}
}

func (p *AuthPlugin) storeLoginHandler() framework.Handler {
	return func(c *framework.Context) error {
		var err error
		loginRequest := &LoginStoreRequest{}

		if err = c.Validate(loginRequest); err != nil {
			return c.WithErrors(err.(validation.Errors)).Back()
		}

		if aUser, err := p.Opts.ResolveUser(loginRequest.Username); aUser != nil {
			log.Println(aUser)
			_, err = p.Login(c.Request().Context(), aUser, loginRequest.Username, loginRequest.Password)
			if err != nil {
				log.Println(err)
				return c.Respond(&framework.R{
					Message:    &framework.AlertMessage{"error", "Login failed."},
					RedirectTo: "/login",
					Payload:    framework.M{"message": "Login failed."},
				})
			}
		} else {
			log.Println(err)
			return c.Respond(&framework.R{
				Message:    &framework.AlertMessage{"error", "Login failed."},
				RedirectTo: "/login",
				Payload:    framework.M{"message": "Login failed."},
			})
		}

		return c.Respond(&framework.R{
			Message:    &framework.AlertMessage{"success", "Login successful."},
			Payload:    framework.M{"message": "Login successful."},
			RedirectTo: p.Opts.HomeRoute,
			Status:     http.StatusOK,
		})
	}
}

func (p *AuthPlugin) Routes() []*framework.Route {
	return []*framework.Route{
		{
			Path:    "/login",
			Method:  "POST",
			Handler: p.Guest(p.storeLoginHandler()),
		},
		{
			Path:    "/login",
			Method:  "GET",
			Handler: p.Guest(p.indexLoginPageHandler()),
		},
		{
			Path:    "/register",
			Method:  "GET",
			Handler: p.Guest(p.indexRegisterPageHandler()),
		},
		{
			Path:    "/register",
			Method:  "POST",
			Handler: p.Guest(p.storeRegisterHandler()),
		},
	}
}

func (p *AuthPlugin) Webhooks() []string {
	return nil
}
