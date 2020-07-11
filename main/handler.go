package main

import (
	"fmt"
	"gin/session"
	"github.com/gin-gonic/gin"
	"net/http"
)
//用户信息
// form 对应html的name
// bingding 是バリデーションチェック
type UserInfo struct {
	Username string `form:"username" binding:"required"`
	Password  string`form:"password" binding:"required"`
}
func indexHandler(ctx *gin.Context){
	ctx.HTML(http.StatusOK, "index.html", nil)
}

func loginHandeler(ctx *gin.Context) {
	if ctx.Request.Method == "POST" {
		var u UserInfo
		err:=ctx.ShouldBind(&u)

		fmt.Println("username: ", u.Username)
		fmt.Println("password:' ", u.Password)
		//sui
		if err != nil {
			fmt.Println("用户名和密码不能为空")
			ctx.HTML(http.StatusOK, "index.html", gin.H{"errMsg":"用户名和密码不能为空"})
			return
		}
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		if username == "yao" && password == "yao" {
			// 从上下文中获取sessionData
			tmpSd, ok := ctx.Get(session.SessionContextName)
			if !ok {
				panic("session middleware")
			}
			// 这句话没有搞明白
			// @TODO
			sd := tmpSd.(session.SessionData)
			// 
			sd.SetKey("isLogin", true)
			sd.Save()
			ctx.Redirect(http.StatusMovedPermanently, ctx.DefaultQuery("next", "index"))
			ctx.HTML(http.StatusOK, "loginSuccess.html", gin.H{"title": "ログイン成功画面", "status": "ログイン成功"})
		} else {
			ctx.HTML(http.StatusOK, "index.html", gin.H{"errMsg":"用户名或者密码不正确"})
		}
	} else {
		ctx.HTML(http.StatusOK, "index.html", nil)
	}

}

func loginCookie01Handler(ctx *gin.Context) {
	cookie, err := ctx.Cookie("key_cookie")
	// no cookie
	if err != nil {
		// name	string	cookie名字
		// value	string	cookie值
		// maxAge	int	有效时间，单位是秒，MaxAge=0 忽略MaxAge属性，MaxAge<0 相当于删除cookie, 通常可以设置-1代表删除，MaxAge>0 多少秒后cookie失效
		// path	string	cookie路径
		// domain	string	cookie作用域
		// secure	bool	Secure=true，那么这个cookie只能用https协议发送给服务器,true的时候只能用https
		// httpOnly	bool	设置HttpOnly=true的cookie不能被js获取到,这样能有效的防止XSS攻击，窃取cookie内容，这样就增加了cookie的安全性
		ctx.SetCookie("key_cookie","value_cookie",
			500, "/", "127.0.0.1", false, true)
	}
	fmt.Println("cookie: ", cookie)
	ctx.HTML(http.StatusOK, "loginCookie.html", gin.H{"title": "ログイン成功画面", "status": "ログイン成功"})
}

func loginCookie02Handler(ctx *gin.Context) {
	var cookieStr *http.Cookie
	if cookie, err := ctx.Request.Cookie("key_cookie02"); err == nil {
		ctx.String(http.StatusOK, cookie.Value)
		cookieStr = cookie
	} else {
		cookie := &http.Cookie{
			Name:  "key_cookie02",
			Value: "value_cookie",
			MaxAge: 48,
			Path: "/",
			Domain: "127.0.0.1",
			Secure: false,
			HttpOnly: true,
		}
		cookieStr = nil
		http.SetCookie(ctx.Writer, cookie)

	}
	fmt.Println("cookie: ", cookieStr)
	ctx.HTML(http.StatusOK, "loginCookie.html", gin.H{"title": "ログイン成功画面", "status": "ログイン成功"})
}

func checkCookieHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "loginCookie.html", gin.H{"title": "ログイン成功画面", "status": "ログイン成功"})
}

// cookie验证
func authCookie() gin.HandlerFunc {
	return func(context *gin.Context) {
		cookie, err := context.Cookie("key_cookie")
		fmt.Println("authCooKIE' ", cookie)
		if err == nil {
			if cookie == "value_cookie" {
				context.Next()
				return
			}
		}
		fmt.Println("authCooKIE' ", cookie)
		context.HTML(http.StatusUnauthorized, "loginCookie.html", gin.H{"title": "ログイン失敗画面", "status": "ログイン失敗"})
		// 验证不通过直接终止请求返回，不做后续的操作
		context.Abort()
		return
	}
}
func vipHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "vip.html", gin.H{"title": "VIP画面", "status": "欢迎VIP用户"})
}

func noRouteHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusNotFound, "error.html", gin.H{"title" : "页面不存在", "errMsg" : "访问的页面不存在"})
}