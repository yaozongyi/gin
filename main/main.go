package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	test()
	router := gin.Default()
	router.LoadHTMLGlob(filepath.Join(os.Getenv("GOPATH"),
		"src/gin/templates/*"))
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/index", indexHandler)
	var username string
	var password string
	router.POST("/login", func(context *gin.Context) {
		username = context.PostForm("username")
		password = context.PostForm("password")
		fmt.Print(username)
		fmt.Println(password)
		context.HTML(http.StatusOK, "loginSuccess.html", gin.H{"title": "ログイン成功画面", "status": "ログイン成功"})
	})

	/**
	* Cookie的生成有两种写法
	 */
	// cookie生成的第一种写法
	// 其实就是第二种写法的封装
	router.GET("/loginCookie01", func(context *gin.Context) {
		cookie, err := context.Cookie("key_cookie")
		// no cookie
		if err != nil {
			// name	string	cookie名字
			// value	string	cookie值
			// maxAge	int	有效时间，单位是秒，MaxAge=0 忽略MaxAge属性，MaxAge<0 相当于删除cookie, 通常可以设置-1代表删除，MaxAge>0 多少秒后cookie失效
			// path	string	cookie路径
			// domain	string	cookie作用域
			// secure	bool	Secure=true，那么这个cookie只能用https协议发送给服务器,true的时候只能用https
			// httpOnly	bool	设置HttpOnly=true的cookie不能被js获取到,这样能有效的防止XSS攻击，窃取cookie内容，这样就增加了cookie的安全性
			context.SetCookie("key_cookie","value_cookie",
				500, "/", "127.0.0.1", false, true)
		}
		fmt.Println("cookie: ", cookie)
		context.HTML(http.StatusOK, "loginCookie.html", gin.H{"title": "ログイン成功画面", "status": "ログイン成功"})
	})
	// cookie生成的第二种写法
	router.GET("/loginCookie02", func(context *gin.Context) {
		var cookieStr *http.Cookie
		if cookie, err := context.Request.Cookie("key_cookie02"); err == nil {
			context.String(http.StatusOK, cookie.Value)
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
			http.SetCookie(context.Writer, cookie)

		}
		fmt.Println("cookie: ", cookieStr)
		context.HTML(http.StatusOK, "loginCookie.html", gin.H{"title": "ログイン成功画面", "status": "ログイン成功"})

	})
	// cookie 验证
	router.GET("/checkCookie", authCookie(), func(context *gin.Context) {
		context.HTML(http.StatusOK, "loginCookie.html", gin.H{"title": "ログイン成功画面", "status": "ログイン成功"})
	})

	/**
	* Session验证
	* 使用 go get github.com/gin-contrib/sessions 插件
	 */
	//
	// 保存到cookie
	// 创建基于cookie的存储引擎，keyPairs1234 参数是用于加密的密钥
	store := cookie.NewStore([]byte("keyPairs1234"))
	// 如果想把session保存到redis把store换成redis就可以了
	//store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	// 设置session
	// name session的name 也是 cookie的name
	// store 是session存储的地方，可以是cookie 也可以是redis
	router.Use(sessions.Sessions("mySession", store))
	router.GET("/loginSession01", func(context *gin.Context) {

		// 初始化session
		session := sessions.Default(context)
		if session.Get("hello") != "word" {
			session.Set("hello", "word")
			session.Save()
		}
		context.JSON(http.StatusOK, gin.H{"hello" : session.Get("hello")})
	})

	// 如果没有匹配到页面则转到404画面
	router.NoRoute(func(context *gin.Context) {
		context.HTML(http.StatusNotFound, "error.html", gin.H{"title" : "页面不存在", "errMsg" : "访问的页面不存在"})
	})
	router.Run()
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