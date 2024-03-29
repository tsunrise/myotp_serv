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

func PrintNewMyError(w http.ResponseWriter, title string, details string, status int) {
	NewMyError(title, details, status).Json(w)
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

func ErrorMissingParam(w http.ResponseWriter, missingParam string) {
	NewMyError("Bad Request", fmt.Sprintf("Missing Parameter %v", missingParam), http.StatusBadRequest).Json(w)
}

func ErrorDatabaseError(w http.ResponseWriter, during string) {
	NewMyError("Database Error", "An error occured when fetching data from the backend database. "+
		"The following process is not successful: "+during, http.StatusInternalServerError).Json(w)
}

func ErrorAuthError(w http.ResponseWriter, text string) {
	NewMyError("Invalid Credentials", text, http.StatusForbidden).Json(w)
}

func ErrorJSONFormatError(w http.ResponseWriter, detail string) {
	PrintNewMyError(w, "JSON Format Error", detail, http.StatusBadRequest)
}

func ErrorRequestMethodError(w http.ResponseWriter, r *http.Request, requiredMethod string) {
	PrintNewMyError(w, "Wrong Request Method",
		fmt.Sprintf("This API only allows %v, got %v", requiredMethod, r.Method),
		http.StatusBadRequest)
}

func ErrorCriticalError(detail string, response http.ResponseWriter) {
	http.Error(response, "Critical Error: "+detail, http.StatusInternalServerError)
}

func ErrorCriticalUnableToWriteResponse(response http.ResponseWriter) {
	ErrorCriticalError("Server is unable to write response.", response)
}
