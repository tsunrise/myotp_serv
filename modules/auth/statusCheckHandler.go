package auth

import (
	"myotp_serv/modules/auth/authUtil"
	"myotp_serv/mydb"
	"myotp_serv/shell"
	"myotp_serv/tokenLib"
	"net/http"
)

func StatusCheckHandler(w http.ResponseWriter, r *http.Request, s *tokenLib.StoreSet, stmt *mydb.StatementsSet) {
	userInfo, err := authUtil.GetUserInfoByRequest(r, s, stmt)
	if err != nil {
		shell.ErrorAuthError(w, err.Error())
		return
	}

	shell.NewResponseStructure(struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Privilege int    `json:"privilege"`
	}{userInfo.UserID, userInfo.UserName, userInfo.Privilege}).Json(w)

}
