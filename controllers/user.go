package controllers

import (
	"beego-demo/models"
	"beego-demo/utils"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type UserController struct {
	BaseController
}

type User struct {
	Id   string `form:"-"`
	Name string `form:"name"`
	Pwd  string `form:"pwd"`
}

/*
  获取用户传递的数据，包括 Get、POST 等方式的请求，均可以使用以下方法获取提交的参数
  GetString(key string) string
  GetStrings(key string) []string
  GetInt(key string) (int64, error)
  GetBool(key string) (bool, error)
  GetFloat(key string) (float64, error)
*/

// @router /user/index [get]
func (c *UserController) Index() {
	c.Data["Method"] = c.Ctx.Request.Method
	c.Data["ClientIP"] = c.getClientIP()
	// 模板传参
	c.Data["Name"] = "denson"
	c.Data["Sex"] = "Male"
	// 设置模板页面
	c.TplName = "user.tpl"
}

// @router /user/post [post]
func (c *UserController) SubmitPost() {
	//var user User
	// 使用c.ParseForm 直接解析到 struct
	//c.ParseForm(&user)
	//fmt.Printf("%v\n", user)

	// 获取POST提交参数
	//  Get、POST 等方式的请求
	id := c.GetString("id")
	name := c.GetString("name")
	pwd := c.GetString("pwd")
	user := &User{
		Id:   id,
		Name: name,
		Pwd:  pwd,
	}
	// 返回JSON数据
	c.Data["json"] = ResponseJson(200, "success", map[string]interface{}{"data": user})
	c.ServeJSON()
}

// @router /user/get [get]
func (c *UserController) SubmitGet() {
	var (
		username string
		age      int
	)
	// 获取GET提交参数
	username = strings.TrimSpace(c.GetString("username"))
	age, _ = c.GetInt("age")
	// Or
	//username := c.GetString("username")
	//age, _ := c.GetInt("age")
	if username == "" {
		c.Data["json"] = ResponseError(4001, "username is required")
		c.ServeJSON()
	}

	if age < 0 || age > 120 {
		c.Data["json"] = ResponseError(4002, "age is invalid！")
		c.ServeJSON()
	}
	// 返回JSON数据
	c.Data["json"] = ResponseJson(200, "success", map[string]interface{}{"username": username, "age": age})
	c.ServeJSON()
}

// @router /register/save [post]
func (c *UserController) SaveRegister() {
	// insert
	var err error
	name := c.GetString("name")
	mobile := c.GetString("mobile")
	password := c.GetString("password")
	fmt.Println(name, mobile, password)
	// 判断
	if mobile == "" {
		c.Data["json"] = ResponseError(4001, "手机号不能为空")
		c.ServeJSON()
	}

	isorno, _ := regexp.MatchString(`^(1[3|4|5|7|8][0-9]\d{4,8})$`, mobile)
	fmt.Println(isorno)
	if !isorno {
		c.Data["json"] = ResponseError(4002, "手机格式本正确")
		c.ServeJSON()
	}

	if password == "" {
		c.Data["json"] = ResponseError(4003, "密码不能为空")
		c.ServeJSON()
	}

	status := models.IsUserMobile(mobile)
	fmt.Println("---status---",status)
	if status {
		c.Data["json"] = ResponseError(4005, "该手机号已经注册过了")
		c.ServeJSON()
	} else {
		err = models.UserSave(name, mobile, MD5V(password))
		if err == nil {
			c.Data["json"] = ResponseJson(200, "添加成功！", nil)
			c.ServeJSON()
		} else {
			c.Data["json"] = ResponseError(5000, err)
			c.ServeJSON()
		}
	}
}

// @router /user/info [get]
func (c *UserController) GetUserInfo()  {
	// query
	id, _ := c.GetInt("id", 1)
	info, err := models.GetUserInfo(id)
	if err != nil {
		c.Data["json"] = ResponseError(400, "")
		c.ServeJSON()
	}
	c.Data["json"] = ResponseJson(200, "success", info)
	c.ServeJSON()
}

// @router /user/list [get]
func (c *UserController) GetUserList()  {
	_, list, err := models.GetUserList()
	if err != nil {
		c.Data["json"] = ResponseError(400, "error")
		c.ServeJSON()
	}
	fmt.Println("xxxxx", utils.GetRandomString(31))
	fmt.Println("xxxxx", utils.TimestampFormat(1595778507, "15:04:05"))
	fmt.Println("xxxxx", utils.UpperFirst("abcdef"))
	c.Data["json"] = ResponseJson(200, "success", list)
	c.ServeJSON()
}

func (c *UserController) DelUser()  {
	//id, _ := c.GetInt("id")
	// strconv.Atoi() 将 string 转化为 int
	// strconv.Itoa() 将 int 转化为 string
	id, err := strconv.Atoi(c.Input().Get("id"))
	//id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.JsonResult(400, "参数错误", "")
	}
	flag := models.DelUserById(id)
	if !flag {
		c.JsonResult(400, "delete error", "")
	}
	c.JsonResult(200, "success", "")
}
