package mydb

import "database/sql"

func initDB(user string, password string, ip string, dbName string) (db sql.DB, err error) {
	panic("to be implemented")
}

type dbError struct {
	text string
}

func newDbError(text string) *dbError {
	return &dbError{text: text}
}

func (e dbError) Error() string {
	return e.text
}


