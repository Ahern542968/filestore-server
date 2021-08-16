package db

import (
	db "filestore-server/db/mysql"
	"fmt"
)

func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	stmt, err := db.DBConn().Prepare(
	"insert ignore into file (`file_sha1`, `file_name`, `file_size`, `file_addr`, `status`) values (?,?,?,?,1)",
	)
	if err != nil {
		fmt.Println("failed to prepare statusment,err:"+err.Error())
		return false
	}
	defer stmt.Close()
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