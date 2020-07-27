package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

type User struct {
	Id       int
	Name     string
	Mobile   string
	Password string
	Status   int
	AddTime  int64
}

type UserInfo struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	AddTime int64  `json:"addTime"`
	Avatar  string `json:"avatar"`
}

func init() {
	// 注册模型
	orm.RegisterModel(new(User))
}

// 根据用户ID获取用户信息
func GetUserInfo(uid int) (UserInfo, error) {
	o := orm.NewOrm()
	var user UserInfo
	err := o.Raw("SELECT id,name,add_time,avatar FROM user WHERE id=? LIMIT 1", uid).QueryRow(&user)
	//obj := &User{
	//	Id: uid,
	//}
	//err = o.Read(obj)
	return user, err
}

func IsUserMobile(mobile string) bool {
	o := orm.NewOrm()
	// 指定需要查询的字段
	user := User{Mobile: mobile}
	err := o.Read(&user, "Mobile")
	fmt.Println("---err---", err)
	if err == orm.ErrNoRows {
		return false
	} else if err == orm.ErrMissPK {
		return false
	}
	return true
}

// 添加用户
func UserSave(name, mobile, password string) error {
	// 获取 ORM 模型
	o := orm.NewOrm()
	// 定义对象
	var user User
	// 给对象赋值
	user.Name = name
	user.Password = password
	user.Mobile = mobile
	user.Status = 1
	user.AddTime = time.Now().Unix()
	// 将对象插入到模型(表), 返回自增id
	_, err := o.Insert(&user)
	return err
}

// 删除用户
func DelUserById(id int) bool {
	// delete
	o := orm.NewOrm()
	var user User
	user.Id = id
	//_, err := o.QueryTable("user").Filter("id", id).Delete()
	if _, err := o.Delete(&user); err != nil {
		return false
	}
	return true
}

// 获取用户列表
func GetUserList() (int64, []UserInfo, error) {
	o := orm.NewOrm()
	// 方法1
	var list []UserInfo
	num, err := o.Raw("select * from user order by id desc").QueryRows(&list)
	return num, list, err
	// 方法2
	//var list []User
	//qs := o.QueryTable("user")
	// qs.Filter("id", id)
	//qs.OrderBy("-id").Limit(pageSize, offset).All(&list)
	//if count, err := qs.All(&list); err != nil {
	//	return count, list, err
	//} else {
	//	return 0, list, err
	//}
}
