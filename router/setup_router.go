package router

import (
	"github.com/gin-gonic/gin"
	"gorm_gin_service/controller"
	"gorm_gin_service/middleware"
)

// SetupApiRouters 用于注册登录的路由
func SetupApiRouters(r *gin.Engine) {
	r.POST("/register", controller.RegisterHandler)
	r.POST("/login", controller.LoginHandler)
	v1 := r.Group("/api/v1")
	r.Use(middleware.AuthMiddleware())
	v1.POST("/project", controller.CreateProjectHandler)
	v1.GET("/project", controller.GetProjectHandler)
	v1.GET("/project/:id", controller.GetProjectDetailHandler)
	v1.PUT("/project", controller.UpdateProjectHandler)
	v1.DELETE("/project/:id", controller.DeleteProjectHandler)
}
