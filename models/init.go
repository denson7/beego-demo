package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func Init()  {
	user := beego.AppConfig.String("user")
	password := beego.AppConfig.String("password")
	host := beego.AppConfig.String("host")
	port := beego.AppConfig.String("port")
	database := beego.AppConfig.String("database")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", user, password, host, port, database)
	fmt.Println("dsn:", dsn)
	// 注册Mysql驱动
	orm.RegisterDriver("mysql", orm.DRMySQL)
	// ORM 必须注册一个名为default的数据库，作为默认使用
	orm.RegisterDataBase("default", "mysql", dsn, 30, 30)
	// 获取 RunMode
	//fmt.Println("RunMode:", beego.BConfig.RunMode)
	//beego.BConfig.RunMode = beego.DEV
	//beego.BConfig.WebConfig.AutoRender = false

	//if beego.BConfig.RunMode == "dev" {
	//	beego.BConfig.WebConfig.DirectoryIndex = true
	//	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	//}
}
