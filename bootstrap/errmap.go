package bootstrap

import (
	"github.com/lemmego/api/app"
	"github.com/lemmego/api/res"
)

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

	if c.WantsJSON() {
		return c.JSON(app.M{"error": msg, "status": status})
	}

	return c.Render(res.NewTemplate(c, page).WithData(map[string]any{
		"title": title, "message": msg,
	}))
}

func LoadErrMap() app.ErrMap {
	return app.ErrMap{
		app.ErrUnauthorized: func(c app.Context) error {
			return errorHandler(c, 401, "401.page.gohtml", "401", "Unauthorized")
		},
		app.ErrForbidden: func(c app.Context) error {
			return errorHandler(c, 403, "403.page.gohtml", "403", "Forbidden")
		},
		app.ErrNotFound: func(c app.Context) error {
			return errorHandler(c, 404, "404.page.gohtml", "404", "Not Found")
		},
		app.ErrPageExpired: func(c app.Context) error {
			return errorHandler(c, 419, "419.page.gohtml", "419", "Page Expired")
		},
		app.ErrInternalServerError: func(c app.Context) error {
			return errorHandler(c, 500, "500.page.gohtml", "500", "Server Error")
		},
	}
}
