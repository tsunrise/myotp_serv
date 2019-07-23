package mydb

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
)

func Install() {
	err := installDB(
		prompt("Enter the Manager/root user name: "),
		prompt("Enter the Manager/root user password: "),
		prompt("Enter MySQL address & port:"),
		prompt("This app will create a database. Enter the database name: "),
		prompt("This app will create a user. Enter the app user name: "),
		prompt("Create the app user password: "),
		prompt("Confirm the password: "))

	//err := installDB("root", "tom100points!", "localhost:3306", "myotp",
	//	"test","testp", "testp")

	if err != nil {
		fmt.Println("❌ An Error Occurred: " + err.Error())
	}

}

func prompt(text string) (ans string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(text)
	scanner.Scan()
	ans = scanner.Text()
	return ans
}

func installDB(managerUser string, managerPassword string, addr string, appDbName string,
	appUser string, appPassword string, appPasswordConfirmed string) error {
	fmt.Println("============================================")
	fmt.Println("🚧 Configuring Database")
	// connect to the database (/)
	db, err := sql.Open("mysql",
		fmt.Sprintf("%v:%v@tcp(%v)/mysql", managerUser, managerPassword, addr))
	if err != nil {
		return err
	}

	// check input format
	if appPassword != appPasswordConfirmed {
		return newDbError("Two passwords mismatch. ")
	}
	if managerUser == "" || managerPassword == "" || addr == "" ||
		appDbName == "" || appUser == "" || appPassword == "" || strings.Contains(appDbName, ";") {
		return newDbError("Invalid Argument")
	}

	// create database
	_, err = db.Exec(fmt.Sprintf("create database %v ;", appDbName))
	if err != nil {
		return newDbError("Fail to create database: " + err.Error())
	}

	// create user
	_, err = db.Exec(fmt.Sprintf("create user '%v'@'%v' identified by '%v'", appUser, "%", appPassword))
	if err != nil {
		return newDbError("Fail to create user: " + err.Error())
	}

	// reconnect
	err = db.Close()
	db2, err := sql.Open("mysql",
		fmt.Sprintf("%v:%v@tcp(%v)/%v", managerUser, managerPassword, addr, appDbName))
	if err != nil {
		return newDbError("Fail to connect to the app database: " + err.Error())
	}

	// create tables
	// users
	_, err = db2.Exec("create table users(" +
		"`user_id` int not null auto_increment,`name` char,`privilege` tinyint default 0,primary key (`user_id`));")
	if err != nil {
		return newDbError("Fail to create users table: " + err.Error())
	}

	// groups
	_, err = db2.Exec("create table `groups` (" +
		"`group_id` int not null auto_increment,`name` char,`user_id` int not null,primary key (`group_id`),foreign key (`user_id`)references users(`user_id`)on delete cascade on update cascade);")
	if err != nil {
		return newDbError("Fail to create groups table: " + err.Error())
	}

	// ticket
	_, err = db2.Exec("create table `ticket` " +
		"( `ticket_index` int not null auto_increment, `id` varchar(512), `token` text, `group_id` int not null, primary key (`ticket_index`), foreign key (`group_id`) references `groups`(`group_id`) on delete cascade on update cascade );")
	if err != nil {
		return newDbError("Fail to create ticket table: " + err.Error())
	}

	_, err = db2.Exec(fmt.Sprintf("grant select, delete, update, insert on %v.* to '%v'@'%v';", appDbName, appUser, "%"))
	if err != nil {
		return newDbError("Fail to grant app user privileges: " + err.Error())
	}

	_, err = db2.Exec("flush privileges;")
	if err != nil {
		return newDbError("Fail to flush privilege: " + err.Error())
	}

	err = db2.Close()
	if err != nil {
		return newDbError("Fail to close connection: " + err.Error())
	}
	fmt.Println("✔ Installation Complete. ")

	info := dbInfo{SqlAddr: addr, DatabaseName:appDbName, AppUserName:appUser, AppUserPassword:appPassword}
	info.writeFile()

	return nil
}

type dbInfo struct {
	SqlAddr string
	DatabaseName string
	AppUserName string
	AppUserPassword string
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
		}else{
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