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
	case urlUtil.MatchExact(path, "ticket/create_group"):
		CreateGroupHandler(w, r, ss, stmt)
	case urlUtil.MatchExact(path, "ticket/view_group"):
		ViewGroupHandler(w, r, ss, stmt)
	case urlUtil.MatchExact(path, "ticket/populate"): // generate tickets
		PopulateTicketsHandler(w, r, ss, stmt, db)
	case urlUtil.MatchExact(path, "ticket/scan"): // verify ticket and num_scanned += 1
		shell.ErrorNotImplemented(w, r, "scan")
	default:
		shell.ErrorNotFound(w, r)
	}
}
