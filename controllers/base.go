package controllers

import (
	"github.com/astaxie/beego"
	"strings"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) JsonResult(code int, msg interface{}, result interface{}) {
	c.Data["json"] = ResponseJson(code, msg, result)
	// 对json进行序列化输出
	c.ServeJSON()
	c.StopRun()
}

// 是否POST提交
func (c *BaseController) IsPost() bool {
	return c.Ctx.Request.Method == "POST"
}

// 获取用户IP地址
func (c *BaseController) getClientIP() string {
	s := strings.Split(c.Ctx.Request.RemoteAddr, ":")
	return s[0]
}
// 获取控制器名称和方法名称
func (c *BaseController) GetControllerNameAndActionName() (string, string) {
	controllerName, actionName := c.GetControllerAndAction()
	return controllerName, actionName
}

// 重定向
func (c *BaseController) redirect(url string) {
	c.Redirect(url, 302)
	c.StopRun()
}