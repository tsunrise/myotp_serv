package main

import (
	"myotp_serv/modules/status"
	"myotp_serv/shell"
	"myotp_serv/util/urlUtil"
	"net/http"
)

func indexRouter(s httpServer, response http.ResponseWriter, request *http.Request) {
	switch path := request.URL.Path; {
	case urlUtil.Match(path, "status"):
		status.Handler(response, request, s.Database)
	default:
		shell.ErrorNotFound(response, request)
	}
}
