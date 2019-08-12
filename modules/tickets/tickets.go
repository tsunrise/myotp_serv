package tickets

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"myotp_serv/modules/auth/authUtil"
	"myotp_serv/mydb"
	"myotp_serv/shell"
	"myotp_serv/tokenLib"
	"myotp_serv/util"
	"net/http"
	"strconv"
	"strings"
)

const maxNumTicketsAtOneTime = 50000
const maxTokenLength = 100
const minTokenLength = 10

func PopulateTicketsHandler(w http.ResponseWriter, r *http.Request,
	s *tokenLib.StoreSet, stmt *mydb.StatementsSet, db *sql.DB) {
	if r.Method != "POST" {
		shell.ErrorRequestMethodError(w, r, "POST")
		return
	}

	userInfo, err := authUtil.GetUserInfoByRequest(r, s, stmt)
	if err != nil {
		shell.ErrorAuthError(w, err.Error())
		return
	}

	postBody := struct {
		GroupID        int    `json:"group_id"`
		TokenLength    int    `json:"token_length"`
		NumTickets     int    `json:"num_tickets"`
		TicketIDFormat string `json:"ticket_id_format"`
	}{}

	// process input json
	err = json.NewDecoder(r.Body).Decode(&postBody)
	if err != nil {
		shell.JSONFormatError(w, err.Error())
		return
	}

	if !strings.Contains(postBody.TicketIDFormat, "*") {
		shell.PrintNewMyError(w, "ticket_id_format is malformed",
			"ticket_id must have * to be replaced by sequence number",
			http.StatusBadRequest)
		return
	}

	if postBody.NumTickets < 0 || postBody.NumTickets > maxNumTicketsAtOneTime {
		shell.PrintNewMyError(w, "Range Error",
			fmt.Sprintf("Number of tickets generated must be in range [0, %v]", maxNumTicketsAtOneTime),
			http.StatusBadRequest)
		return
	}

	if postBody.TokenLength < minTokenLength || postBody.TokenLength > maxTokenLength {
		shell.PrintNewMyError(w, "Range Error",
			fmt.Sprintf("Token Length generated must be in range [%v, %v]", minTokenLength, maxTokenLength),
			http.StatusBadRequest)
		return
	}

	tickets, startID, err := PopulateTicket(userInfo.UserID,
		postBody.GroupID,
		postBody.TokenLength,
		postBody.NumTickets,
		postBody.TicketIDFormat,
		stmt, db, r)

	if err != nil {
		shell.ErrorDatabaseError(w, err.Error())
		return
	}

	shell.NewResponseStructure(struct {
		TicketTokens []string `json:"ticket_tokens"`
		StartID      int      `json:"start_id"`
	}{tickets, startID}).Json(w)

}

func PopulateTicket(uid int, groupID int, tokenLength int, numTicket int, ticketIDFormat string, stmt *mydb.StatementsSet, db *sql.DB, r *http.Request) (tickets []string, startID int, err error) {
	// check if group exist and group correspond to uid
	err = checkValidGID(uid, groupID, stmt)
	if err != nil {
		return
	}

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		return
	}
	insertTicketStmt := tx.Stmt(stmt.InsertTicket)
	tickets = make([]string, 0)
	count, err := getGroupCount(groupID, stmt)
	if err != nil {
		return
	}

	// insert start
	for i := 0; i < numTicket; i++ {
		select {
		// stop transaction when request canceled
		case <-r.Context().Done():
			errRb := tx.Rollback()
			if errRb != nil {
				log.Println("Unable to rollback operation: " + errRb.Error())
			}
			return nil, -1, errors.New("Request canceled. ")
		default:
			token := util.RandBase32Token(tokenLength)
			_, err = insertTicketStmt.Exec(
				strings.Replace(ticketIDFormat, "*", strconv.Itoa(count+i+1), 1), // id
				token, //token
				groupID)
			if err != nil {
				rbErr := tx.Rollback()
				if rbErr != nil {
					err = errors.New(err.Error() + rbErr.Error())
				}
				return
			}
			tickets = append(tickets, token)
		}
	}

	err = tx.Commit()
	startID = count + 1
	if err != nil {
		return
	}

	// return tickets
	return

}

func checkValidGID(uid int, groupID int, stmt *mydb.StatementsSet) (err error) {
	rs, err := stmt.GIDToUID.Query(groupID)
	const errorMessage = "The groupID may not exist or you do not have the right to add tickets. "
	if err != nil {
		return
	}

	if !rs.Next() {
		return errors.New(errorMessage)
	}

	var groupUID int
	err = rs.Scan(&groupUID)
	if err != nil {
		return
	}

	if groupUID != uid {
		return errors.New(errorMessage)
	}

	return
}

func getGroupCount(gid int, stmt *mydb.StatementsSet) (count int, err error) {
	rows, err := stmt.CountGroup.Query(gid)
	if err != nil {
		return
	}
	rows.Next()
	err = rows.Scan(&count)
	if err != nil {
		return
	}
	return

}
