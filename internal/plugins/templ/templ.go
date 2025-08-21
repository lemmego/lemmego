package templ

import (
	"github.com/a-h/templ"
	"github.com/lemmego/api/app"
	"io"
	"net/http"
)

type Templ struct {
	component templ.Component
	ctx       app.Context
}

func New(c app.Context, component templ.Component) *Templ {
	return &Templ{component: component, ctx: c}
}

func (t *Templ) Render(w io.Writer) error {
	t.ctx.SetHeader("content-type", "text/html")
	if t.ctx.Status() == 0 {
		t.ctx.SetStatus(http.StatusOK)
	}
	t.ctx.ResponseWriter().WriteHeader(t.ctx.Status())
	return t.component.Render(t.ctx.RequestContext(), t.ctx.ResponseWriter())
}
