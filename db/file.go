package db

import (
	"database/sql"
	"fmt"

	mydb "filestore-server/db/mysql"
)

func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	stmt, err := mydb.DBConn().Prepare(
	"insert ignore into file (`file_sha1`,`file_name`,`file_size`, `file_addr`,`status`) values (?,?,?,?,1)",
	)
	if err != nil {
		fmt.Println("failed to prepare statement, err:"+err.Error())
		return false
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0{
			fmt.Printf("file with hash: %s has been updated before", filehash)
		}
		return true
	}
	return false
}