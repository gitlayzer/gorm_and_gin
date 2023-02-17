package dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm_gin_service/config"
	"gorm_gin_service/model"
)

func InitMySQL() *gorm.DB {

	// 引用变量并实例化conn
	var conn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		config.DbUser,
		config.DbPass,
		config.DbHost,
		config.DbPort,
		config.DbName,
		config.DbLang,
		config.DbTime,
		config.DbLoc,
	)

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败", err)
	}

	fmt.Println("数据库连接成功")

	if err := db.AutoMigrate(&model.Project{}, &model.User{}); err != nil {
		fmt.Println("表创建失败", err)
	}

	return db
}
