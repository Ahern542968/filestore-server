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


func UserSignIn(username string, password string) bool {
	fmt.Println(username, password)
	stmt, err := mysql.DBConn().Prepare(
		"select * from `user` where `user_name`=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)

	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found: " + username)
		return false
	}

	pRows := mysql.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == password {
		return true
	}
	return false
}

// UpdateToken : 刷新用户登录的token
func UpdateToken(username string, token string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"replace into `user_token` (`user_name`,`user_token`) values (?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}