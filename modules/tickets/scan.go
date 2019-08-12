package tickets

import (
	"encoding/json"
	otp2 "github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"myotp_serv/mydb"
	"myotp_serv/shell"
	"myotp_serv/tokenLib"
	"net/http"
	"time"
)

func ScanHandler(w http.ResponseWriter, r *http.Request, s *tokenLib.StoreSet, stmt *mydb.StatementsSet) {
	// may add scanner authorization

	//use POST
	if r.Method != "POST" {
		shell.ErrorRequestMethodError(w, r, "POST")
		return
	}

	postBody := struct {
		GroupID  int    `json:"group_id"`
		TicketID string `json:"ticket_id"`
		NumPass  string `json:"num_pass"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&postBody)
	if err != nil {
		shell.ErrorJSONFormatError(w, err.Error())
		return
	}

	// get ticket information
	rs, err := stmt.ScanTicket.SelectTicket.Query(postBody.GroupID, postBody.TicketID)
	if err != nil {
		shell.PrintNewMyError(w, "Unable to resolve this ticket", err.Error(), http.StatusInternalServerError)
		return
	}
	if !rs.Next() {
		shell.PrintNewMyError(w, "Ticket Not Exist", "The ticket id does not match any. ", http.StatusForbidden)
		return
	}

	var token string
	var count int
	err = rs.Scan(&token, &count)
	if err != nil {
		shell.ErrorDatabaseError(w, err.Error())
		return
	}

	var resp = struct {
		Valid      bool `json:"valid"`
		NumScanned int  `json:"num_scanned"`
	}{}

	// verify token
	resp.Valid, err = VerifyOTP(token, postBody.NumPass)
	if err != nil {
		shell.NewMyError("Fail to verify otp. ", err.Error(), http.StatusInternalServerError).Json(w)
		return
	}

	if !resp.Valid {
		resp.NumScanned = -1
		shell.NewResponseStructure(resp).Json(w)
		return
	}

	// get number scanned
	resp.NumScanned = count

	// number scanned add plus one
	_, err = stmt.ScanTicket.AddOneScan.Exec(postBody.GroupID, postBody.TicketID)
	if err != nil {
		shell.ErrorDatabaseError(w, "adding number scanned: "+err.Error())
		return
	}

	shell.NewResponseStructure(resp).Json(w)

}

func TokenToOTP(token string) (otp string, err error) {
	return totp.GenerateCodeCustom(token, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp2.DigitsSix,
		Algorithm: otp2.AlgorithmSHA1,
	})
}

func VerifyOTP(token string, otp string) (success bool, err error) {
	success, err = totp.ValidateCustom(otp, token, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp2.DigitsSix,
		Algorithm: otp2.AlgorithmSHA1,
	})

	return
}
