package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm_gin_service/dao"
	"gorm_gin_service/model"
	"net/http"
)

// RegisterHandler 注册
func RegisterHandler(c *gin.Context) {
	p := new(model.User)
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// 使用全局的db对象
	if data := dao.InitMySQL().Create(p); data.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": data.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "注册成功",
		"data": p.Username,
	})
}

// LoginHandler 登录
func LoginHandler(c *gin.Context) {
	// 获取参数并与数据库中的数据进行比对绑定
	p := new(model.User)
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// 定义一个用户对象
	info := &model.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 查询用户是否存在
	if rows := dao.InitMySQL().Where(&info).First(&info); rows == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}

	// 生成token
	token := uuid.New().String()
	// 存入数据库
	if data := dao.InitMySQL().Model(&info).Update("token", token); data.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": data.Error.Error(),
		})
		return
	}

	// 登录成功并返回token
	c.JSON(http.StatusOK, gin.H{
		"msg":  "登录成功",
		"data": token,
	})
}
