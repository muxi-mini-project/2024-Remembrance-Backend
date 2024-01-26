package main

import (
	"flag"
	"remembrance/app/common"
	"remembrance/app/common/config"
	"remembrance/app/common/tube"
	"remembrance/app/core/gorm"
	"remembrance/app/routers"
)

func main() {
	flag.Parse()
	var c config.Config
	//加载图床配置
	tube.Load(c)
	common.DB = gorm.Newgorm() //数据库连接引用全局变量
	e := routers.RouterInit()
	e.Run(":8088")
}
