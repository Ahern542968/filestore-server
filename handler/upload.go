package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
		// 接收文件流及存储到本地目录
		file, header, err:= r.FormFile("file")
		if err != nil {
			fmt.Printf("failed to get upload, err %s\n", err.Error())
			return
		}
		defer file.Close()
		newFile, err := os.Create("./files/"+header.Filename)
		if err != nil {
			fmt.Printf("failed to create file, err %s\n", err.Error())
			return
		}
		defer newFile.Close()
		_, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("failed to save data into file, err %s\n", err.Error())
			return
		}
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload file success")
}
