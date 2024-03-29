package auth

import (
	"myotp_serv/mydb"
	"myotp_serv/shell"
	"myotp_serv/tokenLib"
	"myotp_serv/util/urlUtil"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request, stmt *mydb.StatementsSet, storeSet *tokenLib.StoreSet) {
	switch path := r.URL.Path; {
	case urlUtil.MatchExact(path, "auth/status"):
		StatusCheckHandler(w, r, storeSet, stmt)
	case urlUtil.MatchExact(path, "auth/login"):
		LoginHandler(w, r, stmt, storeSet)
	case urlUtil.MatchExact(path, "auth/create"):
		UserCreateHandler(w, r, stmt)
	default:
		shell.ErrorNotFound(w, r)
	}
}
