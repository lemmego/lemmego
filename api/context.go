package api

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"lemmego/api/db"
	"lemmego/api/fsys"
	"lemmego/api/logger"
	"lemmego/api/vee"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"sync"

	"github.com/golobby/container/v3"
	inertia "github.com/romsar/gonertia"

	"lemmego/api/req"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

func init() {
	gob.Register(&AlertMessage{})
	gob.Register(vee.Errors{})
	gob.Register([]*AlertMessage{})
	gob.Register(vee.Errors{})
}

type Context struct {
	sync.Mutex
	app            *App
	container      container.Container
	request        *http.Request
	responseWriter http.ResponseWriter
}

type AlertMessage struct {
	Type string // success, error, warning, info, debug
	Body string
}

type R struct {
	Status       int
	TemplateName string
	Message      *AlertMessage
	Payload      M
	RedirectTo   string
}

// SetCookie sets a cookie on the response writer
// Example: // c.SetCookie("jwt", token, 60*60*24*7, "/", "", false, true)
func (c *Context) SetCookie(name string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	}
	http.SetCookie(c.responseWriter, cookie)
}

func (c *Context) Alert(typ string, message string) *AlertMessage {
	if typ != "success" && typ != "error" && typ != "warning" && typ != "info" && typ != "debug" {
		return &AlertMessage{Type: "", Body: ""}
	}

	return &AlertMessage{Type: typ, Body: message}
}

func (c *Context) ParseAndValidate(body req.Validator) (any, error) {
	// return error if body is not a pointer
	if reflect.ValueOf(body).Kind() != reflect.Ptr {
		return nil, errors.New("body must be a pointer")
	}

	input, err := c.ParseInput(body)

	if err != nil {
		return nil, err
	}

	if err := req.Validate(c.responseWriter, c.request, input.(req.Validator)); err != nil {
		return nil, err
	}

	return input, nil
}

func (c *Context) ParseInput(inputStruct any) (any, error) {
	input, err := req.ParseInput(c.request, inputStruct)
	if err != nil {
		return inputStruct, err
	}
	return input, nil
}

func (c *Context) Input(inputStruct any) any {
	err := req.In(c, inputStruct)
	if err != nil {
		return nil
	}
	return c.Get(InKey)
}

func (c *Context) SetInput(inputStruct any) error {
	err := req.In(c, inputStruct)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) GetInput() any {
	return c.Get(InKey)
}

func (c *Context) Respond(status int, r *R) error {
	if c.WantsJSON() {
		if r.Payload != nil {
			return c.JSON(r.Status, r.Payload)
		}

		if r.Message.Body != "" {
			return c.JSON(r.Status, M{r.Message.Type: r.Message.Body})
		}
	}

	templateData := &TemplateData{}

	// if r.Message != nil && r.Message.Body != "" {
	// 	c.PutFlash(r.Message.Type, r.Message)
	// }

	if r.RedirectTo != "" {
		var messageType string
		if r.Message != nil {
			messageType = r.Message.Type
		}
		switch messageType {
		case "success":
			c.WithSuccess(r.Message.Body).Redirect(http.StatusFound, r.RedirectTo)
			return nil
		case "info":
			c.WithInfo(r.Message.Body).Redirect(http.StatusFound, r.RedirectTo)
			return nil
		case "warning":
			c.WithWarning(r.Message.Body).Redirect(http.StatusFound, r.RedirectTo)
			return nil
		case "error":
			c.WithError(r.Message.Body).Redirect(http.StatusFound, r.RedirectTo)
			return nil
		default:
			c.Redirect(http.StatusFound, r.RedirectTo)
			return nil
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

func (c *Context) FS() fsys.FS {
	return c.app.fs
}

func (c *Context) DB() *db.DB {
	return c.app.db
}

func (c *Context) I() *inertia.Inertia {
	return c.app.i
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

func (c *Context) WantsJSON() bool {
	return req.WantsJSON(c.request)
}

func (c *Context) JSON(status int, body M) error {
	// TODO: Check if header is already sent
	response, _ := json.Marshal(body)
	c.responseWriter.Header().Set("Content-Type", "application/json")
	c.responseWriter.WriteHeader(status)
	_, err := c.responseWriter.Write(response)
	return err
}

func (c *Context) Send(status int, body []byte) error {
	c.responseWriter.Header().Set("content-type", "text/html")
	c.responseWriter.WriteHeader(status)
	_, err := c.responseWriter.Write(body)
	return err
}

func (c *Context) AuthUser() interface{} {
	return c.Pop("authUser")
}

func (c *Context) resolveTemplateData(data *TemplateData) *TemplateData {
	if data == nil {
		data = &TemplateData{}
	}

	vErrs := vee.Errors{}

	if val, ok := c.Pop("errors").(vee.Errors); ok {
		vErrs = val
	}

	if data.ValidationErrors == nil {
		data.ValidationErrors = vErrs
	}

	data.Messages = append(data.Messages, &AlertMessage{"success", c.PopString("success")})
	data.Messages = append(data.Messages, &AlertMessage{"info", c.PopString("info")})
	data.Messages = append(data.Messages, &AlertMessage{"warning", c.PopString("warning")})
	data.Messages = append(data.Messages, &AlertMessage{"error", c.PopString("error")})

	return data
}

func (c *Context) HTML(status int, body string) error {
	c.responseWriter.Header().Set("Content-Type", "text/html")
	c.responseWriter.WriteHeader(status)
	_, err := c.responseWriter.Write([]byte(body))
	return err
}

func (c *Context) Render(status int, tmplPath string, data *TemplateData) error {
	data = c.resolveTemplateData(data)
	c.responseWriter.Header().Set("Content-Type", "text/html")
	c.responseWriter.WriteHeader(status)
	return RenderTemplate(c.responseWriter, tmplPath, data)
}

func (c *Context) Inertia(status int, filePath string, props map[string]any) error {
	if c.app.i == nil {
		return errors.New("inertia not enabled")
	}
	c.responseWriter.WriteHeader(status)
	return c.App().Inertia().Render(c.ResponseWriter(), c.Request(), filePath, props)
}

func (c *Context) Redirect(status int, url string) {
	c.responseWriter.Header().Set("Location", url)
	c.responseWriter.WriteHeader(status)
}

func (c *Context) WithErrors(errors vee.Errors) *Context {
	if c.app.i != nil {
		c.app.i.ShareProp("errors", errors)
	}
	return c.Put("errors", errors)
}

func (c *Context) WithSuccess(message string) *Context {
	if c.app.i != nil {
		c.app.i.ShareProp("success", message)
	}
	return c.Put("success", message)
}

func (c *Context) WithInfo(message string) *Context {
	if c.app.i != nil {
		c.app.i.ShareProp("info", message)
	}
	return c.Put("info", message)
}

func (c *Context) WithWarning(message string) *Context {
	if c.app.i != nil {
		c.app.i.ShareProp("warning", message)
	}
	return c.Put("warning", message)
}

func (c *Context) WithError(message string) *Context {
	if c.app.i != nil {
		c.app.i.ShareProp("error", message)
	}
	return c.Put("error", message)
}

func (c *Context) WithData(data map[string]any) *Context {
	c.PutFlash("data", data)
	return c
}

func (c *Context) Back(status int) {
	if c.app.i != nil {
		c.App().Inertia().Back(c.ResponseWriter(), c.Request(), status)
		return
	}
	c.Redirect(status, c.request.Referer())
}

func (c *Context) GetParam(key string) string {
	return chi.URLParam(c.request, key)
}

func (c *Context) GetQuery(key string) string {
	return c.request.URL.Query().Get(key)
}

func (c *Context) GetBody() (map[string][]string, error) {
	if err := c.request.ParseForm(); err != nil {
		return nil, err
	}
	return c.request.Form, nil
}

func (c *Context) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	if err := c.request.ParseMultipartForm(32 << 20); err != nil {
		return nil, nil, err
	}
	return c.request.FormFile(key)
}

func (c *Context) Upload(key string, dir string) (*os.File, error) {
	file, header, err := c.FormFile(key)
	defer func() {
		err := file.Close()
		if err != nil {
			logger.V().Info("Form file could not be closed", "Error:", err)
		}
	}()

	if err != nil {
		return nil, err
	}

	return c.FS().Upload(file, header, dir)
}

func (c *Context) Set(key string, value interface{}) {
	c.Lock()
	defer c.Unlock()
	c.request = c.request.WithContext(context.WithValue(c.request.Context(), key, value))
}

func (c *Context) Get(key string) any {
	c.Lock()
	defer c.Unlock()
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
	}
	c.WithError(err.Error()).Back(http.StatusFound)
	return nil
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
	return req.DecodeJSONBody(c.responseWriter, c.request, v)
}
