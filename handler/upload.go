package handler

import (
	"filestore-server/meta"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"filestore-server/util"
)

// UploadHandler ： 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			_, err = io.WriteString(w, "internel server error\n")
			return
		} else {
			_, err = io.WriteString(w, string(data))
		}
	} else if r.Method == "POST" {
		// 接收文件流及存储到本地目
		file, header, err:= r.FormFile("file")
		if err != nil {
			fmt.Printf("failed to get upload, err %s\n", err.Error())
			return
		}
		defer file.Close()

		location := "./files/"+header.Filename

		fileMeta := meta.FileMeta{
			FileName: header.Filename,
			Location: location,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(location)
		if err != nil {
			fmt.Printf("failed to create file, err %s\n", err.Error())
			return
		}
		defer newFile.Close()
		fileMeta.FileSize, err = io.Copy(newFile, file)

		if err != nil {
			fmt.Printf("failed to save data into file, err %s\n", err.Error())
			return
		}
		newFile.Seek(0, 0)

		fileMeta.FileSha1 = util.Sha1(make([]byte, fileMeta.FileSize))

		newFile.Seek(0, 0)

		meta.FileMetas[fileMeta.FileSha1] = fileMeta

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload file success")
}
