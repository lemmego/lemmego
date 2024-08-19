package req

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/romsar/gonertia"

	"github.com/ggicci/httpin"
	"github.com/ggicci/httpin/core"
	"github.com/golang/gddo/httputil/header"
)

const InKey = "input"

type Ctx interface {
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	Set(key string, value interface{})
	Get(key string) interface{}
}

type Validator interface {
	Validate() error
}

type MalformedRequest struct {
	Status  int
	Message string
}

func (mr *MalformedRequest) Error() string {
	return mr.Message
}

func WantsJSON(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "application/json")
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &MalformedRequest{Status: http.StatusUnsupportedMediaType, Message: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &MalformedRequest{Status: http.StatusBadRequest, Message: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &MalformedRequest{Status: http.StatusBadRequest, Message: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %#q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &MalformedRequest{Status: http.StatusBadRequest, Message: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &MalformedRequest{Status: http.StatusBadRequest, Message: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &MalformedRequest{Status: http.StatusBadRequest, Message: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &MalformedRequest{Status: http.StatusRequestEntityTooLarge, Message: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &MalformedRequest{Status: http.StatusBadRequest, Message: msg}
	}

	return nil
}

func HasFormData(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data")
}

func ParseInput(ctx Ctx, inputStruct any, opts ...core.Option) error {
	if !HasFormData(ctx.Request()) && (WantsJSON(ctx.Request()) || gonertia.IsInertiaRequest(ctx.Request())) {
		if err := DecodeJSONBody(ctx.ResponseWriter(), ctx.Request(), inputStruct); err != nil {
			return err
		}
		return nil
	}
	co, err := httpin.New(inputStruct, opts...)

	if err != nil {
		return err
	}

	input, err := co.Decode(ctx.Request())
	if err != nil {
		return err
	}

	reflect.ValueOf(inputStruct).Elem().Set(reflect.ValueOf(input).Elem())

	return nil
}

func In(ctx Ctx, inputStruct any, opts ...core.Option) error {
	if WantsJSON(ctx.Request()) || gonertia.IsInertiaRequest(ctx.Request()) {
		if err := DecodeJSONBody(ctx.ResponseWriter(), ctx.Request(), inputStruct); err != nil {
			return err
		}
		ctx.Set(InKey, inputStruct)
		return nil
	}
	co, err := httpin.New(inputStruct, opts...)

	if err != nil {
		return err
	}

	input, err := co.Decode(ctx.Request())
	if err != nil {
		return err
	}

	ctx.Set(InKey, input)
	return nil
}
