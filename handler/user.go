package handler

import (
	mydb "cloud_storage/db"
	"cloud_storage/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	pwdSalt = "*#890"
)

// 注册接口
func SignupHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.Redirect(http.StatusFound, "/static/view/signup.html")
		return
	}
	c.Request.ParseForm()

	username := c.PostForm("username")
	password := c.PostForm("password")

	if len(username) < 3 || len(password) < 6 {
		c.String(http.StatusOK, "Invalid parameter")
		return
	}

	encPassword := util.Sha1([]byte(password + pwdSalt))

	success := mydb.UserSignup(username, encPassword)
	if success {
		c.String(http.StatusOK, "注册成功")
	} else {
		c.String(http.StatusOK, "注册失败")
	}

}

// 登录接口
func LoginHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "login",
		})
		return
	}
	// 校验
	username := c.PostForm("username")
	password := c.PostForm("password")
	signinRet := mydb.UserSignin(username, password)
	if !signinRet {
		c.String(http.StatusOK, "用户名或密码错误")
		return
	}
	token := GenToken(username)
	upRes := mydb.UpdateToken(username, token)
	if !upRes {
		c.String(http.StatusOK, "登录失败(token生成失败)")
		return
	}

	//3.登录成功后重定向
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + c.Request.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	c.JSON(http.StatusOK, resp.JSONBytes())
}

// UserInfoHandler ： 查询用户信息
func UserInfoHandler(c *gin.Context) {
	username := c.PostForm("username")

	//验证token是否有效
	token := c.PostForm("token")
	isValidToken := IsTokenValid(token)
	if !isValidToken {
		c.String(http.StatusForbidden, "token无效")
		return
	}

	user, err := mydb.GetUserInfo(username)
	if err != nil {
		c.String(http.StatusForbidden, "获取用户信息失败")
		return
	}

	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	c.JSON(http.StatusOK, resp.JSONBytes())
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
