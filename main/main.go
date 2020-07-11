package main

import (
	"gin/session"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob(filepath.Join(os.Getenv("GOPATH"),
		"src/gin/templates/*"))
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	// session作为全局的中间件
	session.InitMgr("memory", "")
	router.Use(session.SessionMiddleware(session.MgrObj))
	router.GET("/index", indexHandler)

	router.Any("/login", loginHandeler)

	/**
	* Cookie的生成有两种写法
	 */
	// cookie生成的第一种写法
	// 其实就是第二种写法的封装
	router.GET("/loginCookie01", loginCookie01Handler)
	// cookie生成的第二种写法
	router.GET("/loginCookie02", loginCookie02Handler)
	// cookie 验证
	router.GET("/checkCookie", authCookie(), checkCookieHandler)
	// vip
	router.GET("/vip", session.AuthMiddleware, vipHandler)

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
	router.NoRoute(noRouteHandler)

	router.Run()
}

