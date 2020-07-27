package utils

import (
	"fmt"
	"strings"
)

// 字符串截取
func SubStr(str interface{}, start, length int) string {
	v := fmt.Sprintf("%v", str)
	if start < 0 {
		start = 0
	}
	slice := strings.Split(v, "")
	l := len(slice)
	if l == 0 || start > l {
		return ""
	}
	if start+length+1 > l {
		return strings.Join(slice[start:], "")
	}
	return strings.Join(slice[start:length], "")
}

// 首字母大写
func UpperFirst(str string) string {
	return strings.Replace(str, str[0:1], strings.ToUpper(str[0:1]), 1)
}
