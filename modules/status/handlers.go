package status

import (
	"database/sql"
	"myotp_serv/shell"
	"myotp_serv/util/urlUtil"
	"net/http"
)

func Handler(response http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch path := r.URL.Path; {
	case urlUtil.MatchExact(path, "status"):
		ans := struct {
			Working bool `json:"working"`
		}{
			true,
		}
		shell.NewResponseStructure(ans).Json(response)
	default:
		shell.ErrorNotFound(response, r)
	}
}
