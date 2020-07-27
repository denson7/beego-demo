package main

import (
	"beego-demo/models"
	_ "beego-demo/routers"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

// 初始化函数
func init() {
	//
}

func main() {
	// 初始化model
	models.Init()
	// 启动服务
	beego.Run()
}
