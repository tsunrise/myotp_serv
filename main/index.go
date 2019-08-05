package main

import (
	"myotp_serv/modules/auth"
	"myotp_serv/modules/status"
	"myotp_serv/modules/tickets"
	"myotp_serv/shell"
	. "myotp_serv/util/urlUtil"
	"net/http"
)

func indexRouter(s httpServer, w http.ResponseWriter, r *http.Request) {
	switch path := r.URL.Path; {
	case Match(path, "status"):
		status.Handler(w, r, s.Database)
	case Match(path, "auth"):
		auth.Handler(w, r, s.DBStatements, s.StoreSet)
	case Match(path, "ticket"):
		tickets.Handler(w, r, s.StoreSet, s.DBStatements, s.Database)
	default:
		shell.ErrorNotFound(w, r)
	}
}
