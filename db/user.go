package db

import (
	"filestore-server/db/mysql"
	"fmt"
)

func UserSignup(username string, password string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into `user` (`user_name`, `user_pwd`, `status`) values (?, ?, 1)",
	)
	if err != nil {
		fmt.Println("failed to insert, error: ", err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println("failed to insert, error: ", err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); err == nil && rowsAffected > 0 {
		return true
	}
	return false
}
