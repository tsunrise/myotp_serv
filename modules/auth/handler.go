package auth

import (
	"myotp_serv/mydb"
	"myotp_serv/shell"
	"myotp_serv/token"
	"myotp_serv/util/urlUtil"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request, stmt *mydb.StatementsSet, storeSet token.StoreSet) {
	switch path := r.URL.Path; {
	case urlUtil.MatchExact(path, "auth/status"):
		shell.ErrorNotImplemented(w, r, "Auth Status Check")
	case urlUtil.MatchExact(path, "auth/login"):
		LoginHandler(w, r, stmt, storeSet)
	default:
		shell.ErrorNotFound(w, r)
	}
}
