package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Response struct {
	w       http.ResponseWriter
	r       *http.Request
	Status  int
	Headers map[string]string
	Body    interface{}
}

type JSONResponse struct {
	http.ResponseWriter
	Status  int
	Headers map[string]string
	Body    map[string]interface{}
}

func NewJSONResponse(w http.ResponseWriter, status int, headers map[string]string, body interface{}) *JSONResponse {
	return &JSONResponse{w, status, headers, body.(map[string]interface{})}
}

func (jr *JSONResponse) Send() {
	response, error := json.Marshal(jr.Body)
	if error != nil {
		log.Println(error)
		jr.Header().Set("Content-Type", "application/json")
		jr.WriteHeader(jr.Status)
		jr.Write([]byte(error.Error()))
		return
	}
	jr.Header().Set("Content-Type", "application/json")
	jr.WriteHeader(jr.Status)
	jr.Write(response)
}

func JSON(w http.ResponseWriter, r *http.Request, status int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// Return http errors
func Error(w http.ResponseWriter, r *http.Request, status int, err error) {
	log.Println(err)
	JSON(w, r, status, map[string]interface{}{"message": err.Error()})
	return
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	message := "Internal Server Error"
	if err != nil && os.Getenv("APP_ENV") != "production" {
		log.Println(err)
		message = err.Error()
	}
	JSON(w, r, http.StatusInternalServerError, map[string]interface{}{"message": message})
	return
}
