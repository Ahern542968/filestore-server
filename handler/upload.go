package handler

import (
	"io"
	"io/ioutil"
	"net/http"
)

// UploadHandler ： 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			_, err = io.WriteString(w, "internel server error")
			return
		} else {
			_, err = io.WriteString(w, string(data))
		}
	} else if r.Method == "POST" {
		// 接收文件流及存储到本地目录
	}
}
