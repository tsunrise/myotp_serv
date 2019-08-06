package auth

import (
	"database/sql"
	"encoding/json"
	"myotp_serv/mydb"
	"myotp_serv/shell"
	"myotp_serv/tokenLib"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request, stmt *mydb.StatementsSet, storeSet *tokenLib.StoreSet) {
	if r.Method != "POST" {
		shell.ErrorRequestMethodError(w, r, "POST")
		return
	}

	postData := struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Hash string `json:"hash"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&postData)
	if err != nil {
		shell.NewMyError("Bad JSON Format", err.Error(), http.StatusBadRequest)
		return
	}

	if postData.ID == "" {
		loginByName(w, r, stmt)
		return
	}

	if postData.Hash == "" {
		shell.ErrorMissingParam(w, "hash")
		return
	}

	rows, err := stmt.CheckUserHashByID.Query(postData.ID, postData.Hash)
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

func saveUserStatus(storeSet *tokenLib.StoreSet, id int) (token string) {
	token = storeSet.Produce()
	userStore, _ := storeSet.Open(token)
	userStore.SetInt("id", id)
	return
}
