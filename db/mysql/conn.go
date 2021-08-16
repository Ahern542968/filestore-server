package mysql

import (
	"fmt"
	"os"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	db, _ := sql.Open("mysql", "root:root@tcp(192.168.10.129:3307)/fileserver?charset=utf8")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("failed to connect to mysql, err:"+err.Error())
		os.Exit(1)
	}
}

func DBConn() *sql.DB {
	return db
}
