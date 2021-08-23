package handler

import (
	"filestore-server/db"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pwdSalt = string("#*2021")
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Redirect(w, r, "/static/view/signup.html", http.StatusFound)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if len(username) < 3 || len(password) < 5 {
		w.WriteHeader(http.StatusBadRequest)
	}

	encPassword := util.Sha1([]byte(password+pwdSalt))

	suc := db.UserSignup(username, encPassword)

	if suc {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}

}


func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传html页面
		data, err := ioutil.ReadFile("/static/view/signin.html")
		if err != nil {
			_, err = io.WriteString(w, "internel server error\n")
		} else {
			_, err = io.WriteString(w, string(data))
		}
		return
	}

	// 1. 校验用户名, 密码
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	encPassword := util.Sha1([]byte(password+pwdSalt))

	pwdChecked := db.UserSignIn(username, encPassword)
	if !pwdChecked {
		w.Write([]byte("FAILED1"))
		return
	}
	// 2. 生成访问凭证

	token := GenToken(username)
	ret := db.UpdateToken(username, token)
	if !ret {
		w.Write([]byte("FAILED2"))
		return
	}

	// 3. 登录成功, 重定向到首页

	respMsg := util.RespMsg{
		Code: 0,
		Msg: "OK",
		Data: struct {
			Username string
			Token string
			Location string
		}{
			Username: username,
			Token: token,
			Location: "http://" + r.Host + "/static/view/home.html",
		},
	}

	w.Write(respMsg.JSONBytes())

}

// GenToken : 生成token
func GenToken(username string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

// IsTokenValid : token是否有效
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	// TODO: 判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	return true
}
