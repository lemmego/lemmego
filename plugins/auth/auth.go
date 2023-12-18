package auth

import (
	"context"
	"errors"
	"fluent-blog/fluent"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"dario.cat/mergo"

	_ "embed"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golobby/container/v3"
	"github.com/invopop/validation"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

var (
	ErrLoginFailed = errors.New("login failed")
	ErrNoStrategy  = errors.New("no strategy provided: either the session manager or the token config must be provided")
	ErrNoSecret    = errors.New("no secret provided: the JWT_SECRET env variable must be provided")
	ErrNoSession   = errors.New("no session provided: the session manager must be provided")
)

//go:embed templates/login.page.tmpl
var loginTmpl []byte

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

type Options struct {
	DB                fluent.DBSession
	Session           *fluent.Session
	TokenConfig       *TokenConfig
	ResolveUser       func(username string, password string) (*AuthUser, error)
	GoogleOAuthConfig *oauth2.Config
	CustomViewMap     map[string]string
	HomeRoute         string
}

type AuthPlugin struct {
	container container.Container
}

func (p *AuthPlugin) SetContainer(c container.Container) {
	p.container = c
}

type Auth struct {
	*AuthPlugin
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

func WithResolveUser(resolveUser func(username string, password string) (*AuthUser, error)) OptFunc {
	return func(opts *Options) {
		opts.ResolveUser = resolveUser
	}
}

func WithSessionManager(session *fluent.Session) OptFunc {
	return func(opts *Options) {
		opts.Session = session
	}
}

func WithTokenConfig(tokenConfig *TokenConfig) OptFunc {
	return func(opts *Options) {
		opts.TokenConfig = tokenConfig
	}
}

func New(opts ...OptFunc) *Auth {
	o := DefaultOptions()

	for _, opt := range opts {
		opt(o)
	}

	// If no ResolveUser function is provided, panic
	if o.ResolveUser == nil {
		panic(errors.New("the ResolveUser function must be provided"))
	}

	// if o.TokenConfig == nil && o.Session == nil {
	// 	panic(ErrNoStrategy)
	// }

	if o.TokenConfig != nil && os.Getenv("JWT_SECRET") == "" {
		panic(ErrNoSecret)
	}

	return &Auth{&AuthPlugin{}, o, nil}
}

func (authn *Auth) Login(ctx context.Context, a *AuthUser, username string, password string) (token string, err error) {
	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)); a.Username != "" && a.Password != "" && err == nil {
		if authn.Opts.Session != nil {
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

func (authn *Auth) ForceLogin(ctx context.Context, a Actor) (token string, err error) {
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

func (authn *Auth) Check(r *http.Request) error {
	if authn.Opts.Session != nil {
		if exists := authn.Opts.Session.Exists(r.Context(), "userId"); exists {
			return nil
		} else {
			return errors.New("user session doesn't exists")
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
		user := &AuthUser{}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			user.ID = claims["sub"].(string)
			user.Username = claims["username"].(string)
		}

		authn.AuthUser = user
	}
	return errors.New("could not parse jwt")
}

func (authn *Auth) Web(next fluent.Handler) fluent.Handler {
	return func(c *fluent.Context) error {
		if err := authn.Check(c.Request()); err != nil {
			if c.WantsJSON() {
				return c.JSON(http.StatusUnauthorized, fluent.M{
					"message": "Unauthenticated.",
				})
			} else {
				return c.Redirect("/login", http.StatusFound)
			}
		} else {
			c.Set("user", authn.AuthUser)
			return next(c)
		}
	}
}

func (authn *Auth) Guest(next fluent.Handler) fluent.Handler {
	return func(c *fluent.Context) error {
		if err := authn.Check(c.Request()); err != nil {
			return next(c)
		} else {
			return c.Back()
		}
	}
}

func (p *Auth) Namespace() string {
	return "fluentcms.auth"
}

func (p *Auth) Init(c container.Container) error {
	p.container = c
	var session *fluent.Session
	c.NamedResolve(session, "session")
	p.Opts.Session = session
	return nil
}

func (p *Auth) EventListeners() map[string]func() {
	return nil
}

func (p *Auth) Migrations() []string {
	return nil
}

func (p *Auth) Templates() map[string][]byte {
	return map[string][]byte{
		"login.page.tmpl": loginTmpl,
	}
}

func (p *Auth) Middlewares() []func(http.Handler) http.Handler {
	return nil
}

func (p *Auth) RouteMiddlewares() map[string]fluent.Middleware {
	return map[string]fluent.Middleware{
		"auth": p.Web,
	}
}

func (p *Auth) showLoginWebHandler() fluent.Handler {
	return func(c *fluent.Context) error {
		vErrs := make(fluent.M)
		if val, ok := c.PopMap("validationErrors").(fluent.M); ok {
			vErrs = val
		}
		err := c.PopString("error")
		return c.Render(200, "login.page.tmpl", &fluent.TemplateData{
			ValidationErrors: vErrs,
			Error:            err,
		})
	}
}

func (p *Auth) loginWebHandler() fluent.Handler {
	return func(c *fluent.Context) error {
		body := c.GetBody()
		loginRequest := &LoginStoreRequest{
			Username: body["username"][0],
			Password: body["password"][0],
		}

		if err := c.Validate(loginRequest); err != nil {
			return c.WithErrors(err.(validation.Errors)).Back()
		}

		aUser, err := p.Opts.ResolveUser(loginRequest.Username, loginRequest.Password)
		_, err = p.Login(c.Request().Context(), aUser, loginRequest.Username, loginRequest.Password)
		if err != nil {
			log.Println(err)
			if c.WantsJSON() {
				return c.JSON(http.StatusInternalServerError, fluent.M{
					"message": "Login failed.",
				})
			} else {
				return c.WithError("Login Failed").Back()
			}
		}

		if c.WantsJSON() {
			return c.JSON(http.StatusOK, fluent.M{
				"message": "login successful",
			})
		} else {
			return c.Redirect(p.Opts.HomeRoute, http.StatusFound)
		}
	}
}

func (p *Auth) Routes() []*fluent.Route {
	return []*fluent.Route{
		{
			Path:    "/login",
			Method:  "POST",
			Handler: p.loginWebHandler(),
		},
		{
			Path:    "/login",
			Method:  "GET",
			Handler: p.showLoginWebHandler(),
		},
	}
}

func (p *Auth) Webhooks() []string {
	return nil
}
