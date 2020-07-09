package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
func indexHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"title": "ログイン画面",
	})
}

 */
func indexHandler(c *gin.Context){
	c.HTML(http.StatusOK, "index.html", nil)
}

func test() {
	fmt.Println("test")
}