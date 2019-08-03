package main

import (
	"myotp_serv/modules/auth"
	"myotp_serv/modules/status"
	"myotp_serv/shell"
	. "myotp_serv/util/urlUtil"
	"net/http"
)

func indexRouter(s httpServer, response http.ResponseWriter, request *http.Request) {
	switch path := request.URL.Path; {
	case Match(path, "status"):
		status.Handler(response, request, s.Database)
	case Match(path, "auth"):
		auth.Handler(response, request, s.DBStatements, s.StoreSet)
	default:
		shell.ErrorNotFound(response, request)
	}
}
