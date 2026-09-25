package bootstrap

import (
	"net/http"

	"github.com/lemmego/api/app"
	"github.com/lemmego/api/res"
	"github.com/lemmego/api/shared"
)

// isInertia reports whether the request came from Inertia, which requires an
// Inertia-shaped response — a page or a redirect. A JSON or HTML error body
// makes the client throw "All Inertia requests must receive a valid Inertia
// response".
func isInertia(c app.Context) bool {
	return c.Header("X-Inertia") != ""
}

func errorHandler(c app.Context, status int, page, title, defaultMsg string) error {
	c.SetStatus(status)

	msg := defaultMsg
	if debug, _ := c.App().Config().Get("app.debug").(bool); debug {
		if e := c.Get("error"); e != nil {
			if he, ok := e.(app.HttpError); ok {
				msg = he.GetHttpMessage().Message
			}
		}
	}

	// A redirect back is a valid Inertia response, and for an expired page it
	// is also the useful one: the form re-renders with the message and the
	// user can simply submit again, which is how Laravel handles 419.
	//
	// The status has to be a 3xx or the client will not follow it, so the
	// error status set above is replaced with 303.
	if isInertia(c) {
		c.PutSession("errors", shared.ValidationErrors{"message": {msg}})
		c.SetStatus(http.StatusSeeOther)
		return c.Back()
	}

	// HTML is checked first. A browser sends
	// "text/html,...,*/*;q=0.8", and the */* in that matches a JSON check, so
	// asking about JSON first served browsers a raw JSON body instead of the
	// error page.
	if c.WantsHTML() {
		return c.Render(res.NewTemplate(c, page).WithData(map[string]any{
			"title": title, "message": msg,
		}))
	}

	return c.JSON(app.M{"error": msg, "status": status})
}

func LoadErrMap() app.ErrMap {
	return app.ErrMap{
		app.ErrUnauthorized:        func(c app.Context) error { return errorHandler(c, 401, "401.page.gohtml", "401", "Unauthorized") },
		app.ErrForbidden:           func(c app.Context) error { return errorHandler(c, 403, "403.page.gohtml", "403", "Forbidden") },
		app.ErrNotFound:            func(c app.Context) error { return errorHandler(c, 404, "404.page.gohtml", "404", "Not Found") },
		app.ErrPageExpired:         func(c app.Context) error { return errorHandler(c, 419, "419.page.gohtml", "419", "Page Expired") },
		app.ErrInternalServerError: func(c app.Context) error { return errorHandler(c, 500, "500.page.gohtml", "500", "Server Error") },
	}
}
