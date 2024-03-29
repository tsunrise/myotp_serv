package mydb

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
)

func InitDB() (*sql.DB, *StatementsSet, error) {
	// establish connection
	info, err := ParseJSONtoDBInfo()
	if err != nil {
		return nil, nil, newDbError(err.Error())
	}
	src := fmt.Sprintf("%v:%v@tcp(%v)/%v", info.AppUserName, info.AppUserPassword, info.SqlAddr, info.DatabaseName)
	db, err := sql.Open("mysql", src)

	if err != nil {
		return nil, nil, newDbError("Fail to connect database: " + err.Error())
	}

	log.Println("Try to connect to database: " +
		fmt.Sprintf("%v@tcp(%v)/%v", info.AppUserName, info.SqlAddr, info.DatabaseName))

	// show table lists
	rows, err := db.Query("show tables")
	if err != nil {
		return nil, nil, newDbError("Fail to access table list: " + err.Error())
	}
	var result string
	for rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			return nil, nil, newDbError("Fail to access table values: " + err.Error())
		}
		log.Println("Database has table: " + result)
	}

	// making statements
	stmts, err := NewStatements(db)
	if err != nil {
		return nil, nil, newDbError("Fail to make prepared statements: " + err.Error())
	}

	return db, stmts, nil

}

func ParseJSONtoDBInfo() (*dbInfo, error) {
	var info *dbInfo
	file, err := ioutil.ReadFile("./db.json")
	if err != nil {
		return nil, NewJSONError(err.Error())
	}
	info, err = jsonToDbInfo([]byte(file))
	if err != nil {
		return nil, NewJSONError(err.Error())
	}

	return info, nil
}

type JSONToDBError struct {
	text string
}

func (e JSONToDBError) Error() string {
	return e.text
}

func NewJSONError(text string) *JSONToDBError {
	return &JSONToDBError{text: "Fail to parse db.json file: " + text}
}
