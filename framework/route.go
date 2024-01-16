package framework

type Route struct {
	Method string
	Path   string
	// Handler     func(http.ResponseWriter, *http.Request)
	Middlewares []Middleware
	Handler     Handler
}

func NewRoute(method string, path string, handler Handler, middlewares ...Middleware) *Route {
	return &Route{
		Method:      method,
		Path:        path,
		Middlewares: middlewares,
		Handler:     handler,
	}
}
