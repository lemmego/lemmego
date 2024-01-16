package framework

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/golobby/container/v3"
	"github.com/invopop/validation"
)

func init() {
	gob.Register(&AlertMessage{})
	gob.Register(&ValidationError{})
	gob.Register([]*AlertMessage{})
	gob.Register([]*ValidationError{})
}

type Context struct {
	app            *App
	container      container.Container
	request        *http.Request
	responseWriter http.ResponseWriter
}

type AlertMessage struct {
	Type string // success, error, warning, info, debug
	Body string
}

type ValidationError struct {
	Key   string
	Value string
}

type R struct {
	Status       int
	TemplateName string
	Message      *AlertMessage
	Payload      M
	RedirectTo   string
}

type ValidationRule struct {
	validation.RuleFunc
}

func (r *ValidationRule) StringEquals(str string) validation.RuleFunc {
	return func(value interface{}) error {
		s, _ := value.(string)
		if s != str {
			return errors.New("unexpected string")
		}
		return nil
	}
}

func (c *Context) Rule() *ValidationRule {
	return &ValidationRule{}
}

func (c *Context) Respond(r *R) error {
	if c.WantsJSON() {
		if r.Payload != nil {
			return c.JSON(r.Status, r.Payload)
		}

		if r.Message.Body != "" {
			return c.JSON(r.Status, M{r.Message.Type: r.Message.Body})
		}
	}

	templateData := &TemplateData{}

	if r.Message != nil && r.Message.Body != "" {
		c.PutFlash(r.Message.Type, r.Message)
	}

	if r.RedirectTo != "" {
		var messageType string
		if r.Message != nil {
			messageType = r.Message.Type
		}
		switch messageType {
		case "success":
			return c.WithSuccess(r.Message.Body).Redirect(r.RedirectTo, http.StatusFound)
		case "info":
			return c.WithInfo(r.Message.Body).Redirect(r.RedirectTo, http.StatusFound)
		case "warning":
			return c.WithWarning(r.Message.Body).Redirect(r.RedirectTo, http.StatusFound)
		case "error":
			return c.WithError(r.Message.Body).Redirect(r.RedirectTo, http.StatusFound)
		default:
			return c.Redirect(r.RedirectTo, http.StatusFound)
		}
	}

	if r.Payload != nil {
		templateData.Data = r.Payload
	}

	if r.TemplateName != "" {
		return c.Render(r.Status, r.TemplateName, &TemplateData{Data: r.Payload})
	}

	return nil
}

func (c *Context) App() *App {
	return c.app
}

func (c *Context) Request() *http.Request {
	return c.request
}

func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

func (c *Context) Templ(component templ.Component) error {
	return component.Render(c.Request().Context(), c.responseWriter)
}

func (c *Context) GetHeader(key string) string {
	return c.request.Header.Get(key)
}

func (c *Context) SetHeader(key string, value string) {
	c.responseWriter.Header().Add(key, value)
}

func (c *Context) Validate(body Validator) error {
	if err := Validate(c.responseWriter, c.request, body); err != nil {
		return err
	}
	return nil
}

func (c *Context) WantsJSON() bool {
	return WantsJSON(c.request)
}

func (c *Context) JSON(status int, body M) error {
	response, _ := json.Marshal(body)
	c.responseWriter.Header().Set("Content-Type", "application/json")
	c.responseWriter.WriteHeader(status)
	c.responseWriter.Write(response)
	return nil
}

func (c *Context) Send(status int, body []byte) error {
	c.responseWriter.Header().Set("content-type", "text/html")
	c.responseWriter.WriteHeader(status)
	_, err := c.responseWriter.Write(body)
	return err
}

func (c *Context) resolveTemplateData(data *TemplateData) *TemplateData {
	if data == nil {
		data = &TemplateData{}
	}
	vErrs := []*ValidationError{}
	if val, ok := c.Pop("validationErrors").([]*ValidationError); ok {
		vErrs = val
	}
	if data.ValidationErrors == nil {
		data.ValidationErrors = []*ValidationError{}
	}

	data.ValidationErrors = append(data.ValidationErrors, vErrs...)

	data.Messages = append(data.Messages, &AlertMessage{"success", c.PopString("success")})
	data.Messages = append(data.Messages, &AlertMessage{"info", c.PopString("info")})
	data.Messages = append(data.Messages, &AlertMessage{"warning", c.PopString("warning")})
	data.Messages = append(data.Messages, &AlertMessage{"error", c.PopString("error")})

	return data
}

func (c *Context) HTML(status int, body string) error {
	c.responseWriter.Header().Set("Content-Type", "text/html")
	c.responseWriter.WriteHeader(status)
	c.responseWriter.Write([]byte(body))
	return nil
}

func (c *Context) Render(status int, tmplPath string, data *TemplateData) error {
	c.responseWriter.Header().Set("Content-Type", "text/html")
	c.responseWriter.WriteHeader(status)
	data = c.resolveTemplateData(data)
	return RenderTemplate(c.responseWriter, tmplPath, data)
}

func (c *Context) Redirect(url string, status int) error {
	c.responseWriter.Header().Set("Location", url)
	c.responseWriter.WriteHeader(status)
	return nil
}

func (c *Context) WithErrors(data map[string]error) *Context {
	errors := []*ValidationError{}
	for key, value := range data {
		errors = append(errors, &ValidationError{key, value.Error()})
	}
	return c.Put("validationErrors", errors)
}

func (c *Context) WithSuccess(message string) *Context {
	c.app.session.Put(c.Request().Context(), "success", message)
	return c
}

func (c *Context) WithInfo(message string) *Context {
	c.app.session.Put(c.Request().Context(), "info", message)
	return c
}

func (c *Context) WithWarning(message string) *Context {
	c.app.session.Put(c.Request().Context(), "warning", message)
	return c
}

func (c *Context) WithError(message string) *Context {
	c.app.session.Put(c.Request().Context(), "error", message)
	return c
}

func (c *Context) WithData(data map[string]error) *Context {
	c.PutFlash("data", data)
	return c
}

func (c *Context) Back() error {
	return c.Redirect(c.GetHeader("Referer"), http.StatusFound)
}

func (c *Context) GetParam(key string) string {
	return chi.URLParam(c.request, key)
}

func (c *Context) GetQuery(key string) string {
	return c.request.URL.Query().Get(key)
}

func (c *Context) GetBody() map[string][]string {
	c.request.ParseForm()
	return c.request.Form
}

func (c *Context) Set(key string, value interface{}) {
	c.request = c.request.WithContext(context.WithValue(c.request.Context(), key, value))
}

func (c *Context) Get(key string) any {
	return c.request.Context().Value(key)
}

func (c *Context) PutFlash(key string, value any) {
	c.app.session.Put(c.Request().Context(), key, value)
}

func (c *Context) Put(key string, value any) *Context {
	c.app.session.Put(c.Request().Context(), key, value)
	return c
}

func (c *Context) Pop(key string) any {
	return c.app.session.Pop(c.Request().Context(), key)
}

func (c *Context) PopString(key string) string {
	return c.app.session.PopString(c.Request().Context(), key)
}

func (c *Context) Error(status int, err error) error {
	if c.WantsJSON() {
		return c.JSON(status, M{"message": err.Error()})
	} else {
		return c.WithError(err.Error()).Back()
	}
}

func (c *Context) InternalServerError(err error) error {
	return c.Error(http.StatusInternalServerError, err)
}

func (c *Context) NotFound(err error) error {
	return c.Error(http.StatusNotFound, err)
}

func (c *Context) BadRequest(err error) error {
	return c.Error(http.StatusBadRequest, err)
}

func (c *Context) Unauthorized(err error) error {
	return c.Error(http.StatusUnauthorized, err)
}

func (c *Context) Forbidden(err error) error {
	return c.Error(http.StatusForbidden, err)
}

func (c *Context) DecodeJSON(v interface{}) error {
	return DecodeJSONBody(c.responseWriter, c.request, v)
}
