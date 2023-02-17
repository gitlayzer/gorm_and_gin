package router

import "github.com/gin-gonic/gin"

// InitRouter 集合所有的路由
func InitRouter() *gin.Engine {

	r := gin.Default()
	
	TestRouter(r)

	SetupApiRouters(r)

	return r
}
