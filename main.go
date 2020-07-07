package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("view/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "ログイン画面",
		})
	})
	var username string
	var password string
	router.POST("/login", func(context *gin.Context) {
		username = context.PostForm("username")
		password = context.PostForm("password")
		fmt.Print(username)
		fmt.Println(password)
		context.HTML(http.StatusOK, "loginSuccess.html", gin.H{"title": "ログイン成功画面", "status": "ログイン成功"})
	})

	router.Run(":8080")
}
