package tickets

import (
	"encoding/json"
	"errors"
	"myotp_serv/modules/auth/authUtil"
	"myotp_serv/mydb"
	"myotp_serv/shell"
	"myotp_serv/tokenLib"
	"net/http"
	"strconv"
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

func ViewGroupHandler(w http.ResponseWriter, r *http.Request, ss *tokenLib.StoreSet, stmt *mydb.StatementsSet) {
	if r.Method != "GET" {
		shell.ErrorRequestMethodError(w, r, "GET")
		return
	}

	userInfo, err := authUtil.GetUserInfoByRequest(r, ss, stmt)
	if err != nil {
		shell.ErrorAuthError(w, err.Error())
		return
	}

	uid := userInfo.UserID

	// support pages
	var page int
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			shell.PrintNewMyError(w, "Bad Page Param", err.Error(), http.StatusBadRequest)
			return
		}

		if page < 1 {
			shell.PrintNewMyError(w, "Bad Page Param", "page >= 1", http.StatusBadRequest)
			return
		}
	} else {
		page = 1
	}

	// process
	groups, err := ViewGroups(uid, page, stmt, r)
	if err != nil {
		shell.PrintNewMyError(w, "Failed to fetch groups", err.Error(), http.StatusInternalServerError)
		return
	}

	shell.NewResponseStructure(struct {
		Groups []group `json:"groups"`
	}{groups}).Json(w)

}

type group struct {
	GroupID   int    `json:"group_id"`
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp"`
}

func ViewGroups(uid int, page int, stmt *mydb.StatementsSet, r *http.Request) (groups []group, err error) {
	offset := (page - 1) * 10
	rows, err := stmt.ViewGroups.Query(uid, offset)
	if err != nil {
		return
	}

	var groupID int
	var time int64
	var name string

	groups = make([]group, 0)
	for rows.Next() {
		select {
		case <-r.Context().Done():
			err = errors.New("Request canceled. ")
			return
		default:
			err = rows.Scan(&groupID, &name, &time)
			if err != nil {
				return
			}
			groups = append(groups, group{groupID, name, time})
		}
	}
	return

}
