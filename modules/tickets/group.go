package tickets

import (
	"encoding/json"
	"myotp_serv/modules/auth/authUtil"
	"myotp_serv/mydb"
	"myotp_serv/shell"
	"myotp_serv/tokenLib"
	"net/http"
)

func CreateGroupHandler(w http.ResponseWriter, r *http.Request, s *tokenLib.StoreSet, stmt *mydb.StatementsSet) {
	if r.Method != "POST" {
		shell.ErrorRequestMethodError(w, r, "POST")
		return
	}

	// get POST data
	info, err := authUtil.GetUserInfoByRequest(r, s, stmt)
	if err != nil {
		shell.ErrorAuthError(w, err.Error())
		return
	}

	postBody := struct {
		GroupName string `json:"group_name"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&postBody)
	if err != nil {
		shell.PrintNewMyError(w, "JSON format error", err.Error(), http.StatusBadRequest)
		return
	}

	if postBody.GroupName == "" {
		shell.ErrorMissingParam(w, "group_name")
		return
	}

	// process
	success, err := CreateGroup(postBody.GroupName, info.UserID, stmt)

	if err != nil {
		shell.PrintNewMyError(w, "Fail to create new group", err.Error(), http.StatusInternalServerError)
		return
	}

	// print
	shell.NewResponseStructure(struct {
		Success bool `json:"success"`
	}{success}).Json(w)

}

func CreateGroup(name string, userID int, stmt *mydb.StatementsSet) (success bool, err error) {
	result, err := stmt.CreateGroup.Exec(name, userID)
	if err != nil {
		return
	}

	n, err := result.RowsAffected()
	if err != nil {
		return
	}

	success = n == 1
	return
}
