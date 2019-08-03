package auth

import (
	"database/sql"
	"myotp_serv/mydb"
	"myotp_serv/shell"
	"myotp_serv/token"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request, stmt *mydb.StatementsSet, storeSet token.StoreSet) {
	id := r.URL.Query().Get("id")
	if id == "" {
		loginByName(w, r, stmt)
		return
	}

	hash := r.URL.Query().Get("hash")
	if hash == "" {
		shell.ErrorMissingParam(w, "hash")
		return
	}

	rows, err := stmt.CheckUserHashByID.Query(id, hash)
	if err != nil {
		shell.ErrorDatabaseError(w, "CheckUserHashByID")
	}

	success, userID, name, privilege, err := getInfo(rows)
	if err != nil {
		shell.ErrorDatabaseError(w, "getUserInfo")
	}

	tokenStr := saveUserStatus(storeSet, userID)

	if success {
		shell.NewResponseStructure(struct {
			Success   bool   `json:"success"`
			UserID    int    `json:"user_id"`
			UserName  string `json:"user_name"`
			Privilege int    `json:"privilege"`
			Token     string `json:"token"`
		}{success, userID, name, privilege, tokenStr}).Json(w)
		return
	}

	shell.NewResponseStructure(struct {
		Success bool `json:"success"`
	}{success}).Json(w)
}

func loginByName(w http.ResponseWriter, r *http.Request, stmt *mydb.StatementsSet) {
	name := r.URL.Query().Get("name")
	if name == "" {
		shell.ErrorMissingParam(w, "id or name")
		return
	}

	hash := r.URL.Query().Get("hash")
	if hash == "" {
		shell.ErrorMissingParam(w, "hash")
		return
	}

	shell.ErrorNotImplemented(w, r, "Login By Name")
}

func getInfo(rows *sql.Rows) (success bool, id int, name string, privilege int, err error) {
	defer rows.Close()
	// check there are matches
	if !rows.Next() {
		success = false
		return
	}

	success = true
	err = rows.Scan(&id, &name, &privilege)
	// now success, id, name, privilege should be ready
	return

}

func saveUserStatus(storeSet token.StoreSet, id int) (token string) {
	token = storeSet.Produce()
	userStore, _ := storeSet.Open(token)
	userStore.SetInt("id", id)
	return
}
