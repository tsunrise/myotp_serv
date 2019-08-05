package tickets

import (
	"database/sql"
	"myotp_serv/mydb"
	"myotp_serv/shell"
	"myotp_serv/tokenLib"
	"myotp_serv/util/urlUtil"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request, ss *tokenLib.StoreSet, stmt *mydb.StatementsSet, db *sql.DB) {
	switch path := r.URL.Path; {
	case urlUtil.MatchExact(path, "ticket/new_group"):
		shell.ErrorNotImplemented(w, r, "new_group")
	case urlUtil.MatchExact(path, "ticket/populate"): // generate tickets
		shell.ErrorNotImplemented(w, r, "populate")
	case urlUtil.MatchExact(path, "ticket/scan"): // verify ticket and num_scanned += 1
		shell.ErrorNotImplemented(w, r, "scan")
	default:
		shell.ErrorNotFound(w, r)
	}
}
