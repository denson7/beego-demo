package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"time"
)

// 定义返回数据 json 格式
type JsonStruct struct {
	Code   int         `json:"code"`
	Msg    interface{} `json:"msg"`
	Result interface{} `json:"result"`
}

// 返回 json 数据
func ResponseJson(code int, msg interface{}, result interface{}) (json *JsonStruct) {
	json = &JsonStruct{Code: code, Msg: msg, Result: result}
	return
}

func ResponseError(code int, msg interface{}) (json *JsonStruct) {
	json = &JsonStruct{Code: code, Msg: msg}
	return
}

func MD5V(password string) string {
	h := md5.New()
	h.Write([]byte(password + beego.AppConfig.String("md5code")))
	return hex.EncodeToString(h.Sum(nil))
}

//格式化时间
func DateFormat(times int64) string {
	videoTime := time.Unix(times, 0)
	return videoTime.Format("2006-01-02")
}

