package mydb

import (
	"encoding/json"
	"fmt"
	"os"
)

type dbInfo struct {
	SqlAddr         string `json:"sql_addr"`
	DatabaseName    string `json:"database_name"`
	AppUserName     string `json:"app_user_name"`
	AppUserPassword string `json:"app_user_password"`
}

func jsonToDbInfo(b []byte) (*dbInfo, error) {
	var info dbInfo
	err := json.Unmarshal(b, &info)
	if err != nil {
		return nil, err
	}
	return &info, err

}

func (d dbInfo) json() []byte {
	ans, _ := json.Marshal(d)
	return ans
}

func (d dbInfo) writeFile() {
	data := d.json()
	f, err := os.Create("./db.json")
	if err != nil {
		fmt.Println("⚠ Warning: Unable to create db.json. ")
		return
	}
	// defer closing file
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("⚠ Unable to close file stream.")
		} else {
			fmt.Println("✔ Database information has been saved to db.json.")
		}
	}()

	_, err = f.Write(data)
	if err != nil {
		fmt.Println("⚠ Warning: Unable to write db.json. ")
		return
	}

	err = f.Sync()
	if err != nil {
		fmt.Println("⚠ Warning: Fail to save file db.json. " + err.Error())
		return
	}

}
