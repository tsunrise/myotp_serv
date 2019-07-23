package test

import (
	"myotp_serv/mydb"
	"testing"
	"time"
)

func TestParseJson(t *testing.T) {
	info, err := mydb.ParseJSONtoDBInfo()
	if err != nil {
		t.Error(err.Error())
		t.Fail()
		return
	}
	if info.SqlAddr != "localhost:3306" {
		t.Fail()
	}
}

func TestStmt(t *testing.T) {
	db, stmts, err := mydb.InitDB()
	if err != nil {
		t.Error(err.Error())
		t.Fail()
		return
	}
	go func() {
		_, err = stmts.NewUser.Exec("jerry1")
		if err != nil {
			t.Error(err.Error())
			t.Fail()
			return
		}
	}()
	go func() {
		_, err = stmts.NewUser.Exec("jerry2")
		if err != nil {
			t.Error(err.Error())
			t.Fail()
			return
		}
	}()
	go func() {
		_, err = stmts.NewUser.Exec("jerry3")
		if err != nil {
			t.Error(err.Error())
			t.Fail()
			return
		}
	}()

	time.Sleep(3 * time.Second)
	_, _ = db.Exec("delete from users where name like 'jerry%';")
}
