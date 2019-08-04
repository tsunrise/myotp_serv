package authUtil

import (
	"errors"
	"myotp_serv/mydb"
	"myotp_serv/tokenLib"
	"net/http"
)

type UserInfo struct {
	UserID    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	Privilege int    `json:"privilege"`
	Hash      string `json:"hash"`
}

func GetUserInfoByToken(token string, s *tokenLib.StoreSet, stmt *mydb.StatementsSet) (info UserInfo, err error) {
	userStore, err := s.Open(token)
	if err != nil {
		return
	}

	id, ok := userStore.GetInt("id")
	if !ok {
		err = errors.New("This token appears to be broken: user id information does not exist. ")
		return
	}

	var name, hash string
	var privilege int
	rows, err := stmt.SelectUser.Query(id)
	if err != nil {
		return
	}

	if !rows.Next() {
		err = errors.New("This token appears to be broken: user id information is invalid. ")
		return
	}

	err = rows.Scan(&name, &privilege, &hash)

	if err != nil {
		return
	}

	info = UserInfo{
		id, name, privilege, hash,
	}
	return
}

func GetUserInfoByRequest(r *http.Request, s *tokenLib.StoreSet, stmt *mydb.StatementsSet) (info UserInfo, err error) {
	token := r.URL.Query().Get("token")
	if token == "" {
		err = errors.New("missing parameter: token")
		return
	}

	info, err = GetUserInfoByToken(token, s, stmt)
	return

}
