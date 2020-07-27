package routers

import (
	"beego-demo/controllers"
	"github.com/astaxie/beego"
)

func init() {
	// 加载路由
    //beego.Router("/", &controllers.MainController{})
    // 固定路由:即全匹配的路由
	beego.Router("/user/index", &controllers.UserController{}, "get:Index")
	beego.Router("/user/get", &controllers.UserController{}, "get:SubmitGet")
	beego.Router("/user/post", &controllers.UserController{}, "post:SubmitPost")
	beego.Router("/user/info", &controllers.UserController{}, "get:GetUserInfo")
	beego.Router("/user/list", &controllers.UserController{}, "get:GetUserList")
	beego.Router("/user/del", &controllers.UserController{}, "get:DelUser")
	beego.Router("/register/save", &controllers.UserController{}, "post:SaveRegister")
	beego.Router("/api/upload", &controllers.UploadController{}, "post:Upload")
    // 基本路由
	//beego.Get("/get",func(ctx *context.Context){
	//	ctx.Output.Body([]byte("hello get"))
	//})
	//beego.Post("/post",func(ctx *context.Context){
	//	ctx.Output.Body([]byte("hello post"))
	//})
    // 正则路由
	//beego.Router("/api/?:id", &controllers.UserController{})

	// 注解路由
	//beego.Include(&controllers.UserController{})
	// 命名空间
	//ns := beego.NewNamespace("/v1",
	//	beego.NSNamespace("/user",
	//		beego.NSInclude(
	//			&controllers.UserController{},
	//		),
	//      beego.NSRouter("/login", &controllers.UserController{}, "post:Login"),
	//	))
	//beego.AddNamespace(ns)
}
