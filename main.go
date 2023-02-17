package main

import (
	"gorm_gin_service/dao"
	"gorm_gin_service/router"
)

func main() {
	dao.InitMySQL()
	r := router.InitRouter()
	r.Run(":80")
}
