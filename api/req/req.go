package req

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

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
	Status int
	Msg    string
}

func (mr *MalformedRequest) Error() string {
	return mr.Msg
}

func WantsJSON(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "application/json")
}

func Validate(w http.ResponseWriter, r *http.Request, body Validator) error {
	if WantsJSON(r) {
		if err := DecodeJSONBody(w, r, body); err != nil {
			return err
		}
	}
	if err := body.Validate(); err != nil {
		return err
	}
	return nil
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &MalformedRequest{Status: http.StatusUnsupportedMediaType, Msg: msg}
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
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %#q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &MalformedRequest{Status: http.StatusRequestEntityTooLarge, Msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	return nil
}

func ParseInput[T any](r *http.Request, inputStruct T, opts ...core.Option) (T, error) {
	co, err := httpin.New(inputStruct, opts...)

	if err != nil {
		return inputStruct, err
	}

	input, err := co.Decode(r)
	if err != nil {
		return inputStruct, err
	}

	return input.(T), nil
}

func In(ctx Ctx, inputStruct any, opts ...core.Option) error {
	if WantsJSON(ctx.Request()) {
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