package mydb

import "database/sql"

type StatementsSet struct {
	NewUser     *sql.Stmt //name, hash
	DeleteUser  *sql.Stmt //user_id
	PromoteUser *sql.Stmt //user_id (set privilege to 1)
	FireUser    *sql.Stmt //user_id (set privilege to 0)
}

func NewStatements(db *sql.DB) (*StatementsSet, error) {
	newUser, err := db.Prepare("insert into users(`name`, `hash`) values (?, ?);")
	if err != nil {
		return nil, err
	}
	deleteUser, err := db.Prepare("delete from users where user_id=?")
	if err != nil {
		return nil, err
	}
	promoteUser, err := db.Prepare("update users set privilege=1 where user_id=?")
	if err != nil {
		return nil, err
	}
	fireUser, err := db.Prepare("update users set privilege=0 where user_id=?")
	if err != nil {
		return nil, err
	}

	return &StatementsSet{
		NewUser:     newUser,
		DeleteUser:  deleteUser,
		PromoteUser: promoteUser,
		FireUser:    fireUser,
	}, nil
}
