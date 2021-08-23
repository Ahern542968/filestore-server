package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"filestore-server/meta"
	"filestore-server/util"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传html页面
		data, err := ioutil.ReadFile("/static/view/index.html")
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
		fileMeta.FileSha1 = util.FileSha1(newFile)

		//meta.FileMetas[fileMeta.FileSha1] = fileMeta
		_ = meta.UpdateFileMetaDB(fileMeta)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload file success")
}


func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileHash := r.Form["fileHash"][0]
	fileMeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if fileMeta != nil {
		data, err := json.Marshal(fileMeta)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	} else {
		w.Write([]byte(`{"code":-1,"msg":"no such file"}`))
	}
}


func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileHash := r.Form.Get("fileHash")
	fileMeta := meta.GetFileMeta(fileHash)
	f, err := os.Open(fileMeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fileMeta.FileName+"\"")
	w.Write(data)
}


func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileHash := r.Form.Get("fileHash")
	opType := r.Form.Get("op")
	newFileName := r.Form.Get("filename")
	if opType != "0" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	currentFileMeta := meta.GetFileMeta(fileHash)
	currentFileMeta.FileName = newFileName
	meta.UpdateFileMeta(currentFileMeta)

	data, err := json.Marshal(currentFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileHash := r.Form.Get("fileHash")

	fileMeta := meta.GetFileMeta(fileHash)
	os.Remove(fileMeta.Location)
	meta.RemoveFileMeta(fileHash)
	w.WriteHeader(http.StatusOK)

}