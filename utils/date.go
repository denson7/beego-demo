package utils

import "time"

// 时间戳格式化
//@timestamp  int      时间戳 秒
//@format     string   时间格式，如："2006-01-02 15:04:05"
func TimestampFormat(timestamp int, format ...string) string {
	formats := "2006-01-02 15:04:05"
	if len(format) > 0 {
		formats = format[0]
	}
	return time.Unix(int64(timestamp), 0).Format(formats)
}
// 获取Y-M-D日期格式
func GetDateYMD(timestamp int64) string {
	if timestamp <= 0 {
		return ""
	}
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02")
}