package session

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SessionCookieName = "session_cookie_id"
	SessionContextName = "session_context_name"
)

type Mgr interface {
	Init(addr string, option ...string)
	GetSessionData(sessionId string) (sd SessionData, err error)
	createSession()(sd SessionData)
}
type SessionData interface {
	// 获取ID
	GetId() string
	// 获取value
	GetKey(key string) (value interface{}, err error)
	// 设置 key value
	SetKey(key string, value interface{})
	// 删除
	DelKey(key string)
	// 保存
	Save()
}
var (
	MgrObj Mgr
)

// 构造一个Mgr
func InitMgr(name string, addr string, option ...string) {
	switch name {
	case "memory":
		MgrObj = NewMemory()
	}
	MgrObj.Init(addr, option ...)
}

/**
* Session
 */
func SessionMiddleware(mgrObj Mgr) gin.HandlerFunc {
	return func(context *gin.Context) {
		// 请求过来从请求的cookie中获取sessionId
		SessionID, err := context.Cookie(SessionCookieName)
		var sd SessionData
		if err != nil {
			// 第一次发送请求，给用户建一个SessionData，分配一个SessionId
			sd = mgrObj.createSession()
			//更新sessionid
			SessionID=sd.GetId()//这个sessionid用于回写coookie
		} else {
			//取到SessionId
			//根据SessionId去仓库取SessionData
			sd , err = mgrObj.GetSessionData(SessionID)
			if err != nil {
				//SessionId有误，取不到SessionData,可能是自己伪造的
				//重建SessionData
				sd = mgrObj.createSession()
				//更新sessionid
				SessionID=sd.GetId()//这个sessionid用于回写coookie
			}
		}
		// 如何实现让后续所有请求的方法都拿到sessiondata？让每个用户的sessiondata都不同
		// 利用gin框架的c.Set("session", sessionData)
		context.Set(SessionContextName, sd)
		
		// 回写到Cookie
		context.SetCookie(SessionCookieName, SessionID, 3600, "/", "127.0.0.1", false, true)
		context.Next()
	}
}
/**
* 校验是否登录
 */
func AuthMiddleware(context *gin.Context){
	// 从上下文中获取sessionData
	tmpSd, ok := context.Get(SessionContextName)
	if !ok {
		panic("session middleware")
	}
	// 这句话没有搞明白
	// @TODO
	sd := tmpSd.(SessionData)
		//fmt.Printf("%v\n", sd)
		value, err := sd.GetKey("isLogin")
		if err != nil {
			fmt.Println(err)
			context.Redirect(http.StatusFound, "/login")
			return
		}
		fmt.Println(value)
		isLogin, ok := value.(bool)
		if !ok {
			context.Redirect(http.StatusFound, "/login")
			return
		}
		if !isLogin {
			context.Redirect(http.StatusFound, "/login")
			return
		}
		context.Next()
}