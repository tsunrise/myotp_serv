package auth

import (
	"encoding/json"
	"myotp_serv/mydb"
	"myotp_serv/shell"
	"net/http"
)

func UserCreateHandler(w http.ResponseWriter, r *http.Request, stmt *mydb.StatementsSet) {
	switch r.Method {
	case "POST":
		gotData := struct {
			Name string `json:"name"`
			Hash string `json:"hash"`
		}{}
		err := json.NewDecoder(r.Body).Decode(&gotData)
		if err != nil {
			shell.PrintNewMyError(w, "Bad Json Structure", err.Error(), http.StatusBadRequest)
			return
		}

		name, hash := gotData.Name, gotData.Hash
		if name == "" || hash == "" {
			shell.ErrorMissingParam(w, "json: name or hash")
			return
		}

		success, err := CreateUser(name, hash, stmt)
		if err != nil {
			shell.PrintNewMyError(w, "Fail to create new user", "This username may already exist. Detail: "+err.Error(), http.StatusInternalServerError)
			return
		}

		shell.NewResponseStructure(struct {
			Success bool `json:"success"`
		}{success}).Json(w)
	default:
		shell.ErrorRequestMethodError(w, r, "POST")

	}
}

func CreateUser(name string, hash string, stmt *mydb.StatementsSet) (success bool, err error) {
	result, err := stmt.NewUser.Exec(name, hash)
	if err != nil {
		return
	}
	n, _ := result.RowsAffected()
	success = n > 0

	return
}
