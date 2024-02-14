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
	//读取配置文件
	config.Loadconfig()
	//fmt.Println(common.CONFIG)
	flag.Parse()
	//加载图床配置
	tube.Load()
	common.RDB = gorm.Newredis() //连接redis
	common.DB = gorm.Newmysql()  //连接mysql引用全局变量
	e := routers.RouterInit()
	e.Run(":8088")
}
