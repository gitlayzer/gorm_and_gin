package controller

import (
	"github.com/gin-gonic/gin"
	"gorm_gin_service/dao"
	"gorm_gin_service/model"
	"net/http"
	"strconv"
)

// CreateProjectHandler 新增项目
func CreateProjectHandler(c *gin.Context) {
	p := new(model.Project)
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// 数据落库
	if data := dao.InitMySQL().Create(&p); data.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": data.Error.Error(),
		})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{
		"msg":  "项目创建成功",
		"data": p,
	})
}

// GetProjectHandler 查询项目
func GetProjectHandler(c *gin.Context) {
	projects := make([]model.Project, 0)
	// 列出所有项目
	if err := dao.InitMySQL().Find(&projects); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error.Error(),
		})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{
		"msg":  "项目列表查询成功",
		"data": projects,
	})
}

// GetProjectDetailHandler 查看项目详情
func GetProjectDetailHandler(c *gin.Context) {
	projectId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	p := new(model.Project)
	// 查看项目详情
	if err := dao.InitMySQL().Where("id = ?", projectId).First(&p); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error.Error(),
		})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{
		"msg":  "项目信息查询成功",
		"data": p,
	})
}

// UpdateProjectHandler 更新项目
func UpdateProjectHandler(c *gin.Context) {
	p := new(model.Project)
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// 数据落库
	if err := dao.InitMySQL().Save(&p); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error.Error(),
		})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{
		"msg": "项目更新成功",
	})
}

// DeleteProjectHandler 删除项目
func DeleteProjectHandler(c *gin.Context) {
	projectId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// 删除关联表的数据
	if err := dao.InitMySQL().Select("Users").Delete(&model.Project{ID: projectId}); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error.Error(),
		})
		return
	}

	// 数据库删除
	if err := dao.InitMySQL().Delete(&model.Project{}, projectId); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error.Error(),
		})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{
		"msg": "项目删除成功",
	})
}
