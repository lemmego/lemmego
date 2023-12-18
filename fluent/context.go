package fluent

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"
)

type Context struct {
	app            *App
	request        *http.Request
	responseWriter http.ResponseWriter
}

func (c *Context) App() *App {
	return c.app
}

func (c *Context) Request() *http.Request {
	return c.request
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

func (c *Context) Render(status int, tmplPath string, data *TemplateData) error {
	return RenderTemplate(c.responseWriter, tmplPath, data)
}

func (c *Context) Redirect(url string, status int) error {
	c.responseWriter.Header().Set("Location", url)
	c.responseWriter.WriteHeader(status)
	return nil
}

func (c *Context) WithErrors(data map[string]error) *Context {
	flashable := make(M)
	for key, value := range data {
		flashable[key] = value.Error()
	}
	c.PutFlash("validationErrors", flashable)
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

// func (c *Context) GetBodyJSON(v interface{}) error {
// 	if err := json.Unmarshal(c.GetBody(), v); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *Context) GetBodyStruct(v interface{}) error {
// 	if err := json.Unmarshal(c.GetBody(), v); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *Context) GetBodyMap() (M, error) {
// 	var body M
// 	if err := json.Unmarshal(c.GetBody(), &body); err != nil {
// 		return nil, err
// 	}
// 	return body, nil
// }

func (c *Context) Set(key string, value interface{}) {
	c.request = c.request.WithContext(context.WithValue(c.request.Context(), key, value))
}

func (c *Context) Get(key string) any {
	return c.request.Context().Value(key)
}

func (c *Context) PutFlash(key string, value any) {
	if reflect.TypeOf(value).Kind() == reflect.Struct || reflect.TypeOf(value).Kind() == reflect.Map {
		json, _ := json.Marshal(value)
		value = string(json)
	}

	c.app.session.Put(c.Request().Context(), key, value)
}

func (c *Context) PutMap(key string, value map[string]any) *Context {
	json, _ := json.Marshal(value)
	c.app.session.Put(c.Request().Context(), key, string(json))
	return c
}

func (c *Context) PopMap(key string) any {
	val := c.app.session.Pop(c.Request().Context(), key)
	if val == nil {
		return nil
	}
	if reflect.TypeOf(val).Kind() == reflect.String {
		var obj M
		json.Unmarshal([]byte(val.(string)), &obj)
		return obj
	}

	return val
}

func (c *Context) PopString(key string) string {
	str := c.app.session.Pop(c.Request().Context(), key)
	if val, ok := str.(string); ok {
		return val
	}
	return ""
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
