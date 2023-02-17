package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm_gin_service/dao"
	"gorm_gin_service/model"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		selectToken := &model.User{}
		row := dao.InitMySQL().Where("token = ?", token).First(selectToken).RowsAffected
		if row == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "Token is invalid",
			})
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}
