package test

import (
	"myotp_serv/mydb"
	"testing"
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
	process := make(chan *interface{})
	go func() {
		_, err = stmts.NewUser.Exec("jerry1", "sgfdljhk")
		if err != nil {
			t.Error(err.Error())
			t.Fail()
			return
		}
		process <- nil
	}()
	go func() {
		_, err = stmts.NewUser.Exec("jerry2", "adflhkjs")
		if err != nil {
			t.Error(err.Error())
			t.Fail()
			return
		}
		process <- nil
	}()
	go func() {
		_, err = stmts.NewUser.Exec("jerry3", "wtprueio")
		if err != nil {
			t.Error(err.Error())
			t.Fail()
			return
		}
		process <- nil
	}()

	for i := 0; i < 3; i++ {
		<-process
	}
	_, err = db.Exec("delete from users where name like 'jerry%';")
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	}
}
