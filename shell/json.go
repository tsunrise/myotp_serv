package shell

import (
	"encoding/json"
	"net/http"
	"time"
)

type ResponseStructure struct {
	Error     bool
	Timestamp int64
	Body      interface{}
}

func NewResponseStructure(body interface{}) *ResponseStructure {
	return &ResponseStructure{Error: false, Timestamp: time.Now().Unix(), Body: body}
}



func (r ResponseStructure) Json(response http.ResponseWriter) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	err := json.NewEncoder(response).Encode(r)
	if err != nil {
		ErrorCriticalError("Unable to produce json-format response.", response)
		return
	}
}