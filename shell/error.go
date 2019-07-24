package shell

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MyError struct {
	Error     bool
	Status    int
	Title     string
	Details   string
	Timestamp int64
}

func NewMyError(title string, details string, status int) *MyError {
	return &MyError{Error: true, Status: status, Title: title, Details: details, Timestamp: time.Now().Unix()}
}

func (e MyError) Json(response http.ResponseWriter) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(e.Status)
	err := json.NewEncoder(response).Encode(e)
	if err != nil {
		ErrorCriticalError("Unable to produce json-format Error.", response)
		return
	}
}

func ErrorNotFound(response http.ResponseWriter, request *http.Request) {
	e := NewMyError("Not Found",
		fmt.Sprintf("The API call (%v) is not listed in the documentation. ", request.URL.Path),
		http.StatusNotFound)
	e.Json(response)

}

func ErrorNotImplemented(w http.ResponseWriter, r *http.Request, serviceName string) {
	e := NewMyError("Service Not Available",
		fmt.Sprintf("Service \"%v\" is under development. ", serviceName), http.StatusServiceUnavailable)
	e.Json(w)
}

func ErrorCriticalError(detail string, response http.ResponseWriter) {
	http.Error(response, "Critical Error: "+detail, http.StatusInternalServerError)
}

func ErrorCriticalUnableToWriteResponse(response http.ResponseWriter) {
	ErrorCriticalError("Server is unable to write response.", response)
}
