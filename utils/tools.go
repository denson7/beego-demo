package utils

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//生成随机字符串
func GetRandomString(lens int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lens; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// 转换字节大小
// size 单位为 B
func FormatByte(size int) string {
	fSize := float64(size)
	// 字节单位
	units := [6]string{"B", "KB", "MB", "GB", "TB", "PB"}
	var i int
	for i = 0; fSize >= 1024 && i < 5; i++ {
		fSize /= 1024
	}
	num := fmt.Sprintf("%.2f", fSize)
	return string(num) + " " + units[i]
}

//获取允许上传文件的类型
func GetUploadFileExt() []string {
	ext := beego.AppConfig.DefaultString("upload_file_ext", "png|jpg|jpeg|gif|txt|doc|docx|pdf")
	temp := strings.Split(ext, "|")
	exts := make([]string, len(temp))

	i := 0
	for _, item := range temp {
		if item != "" {
			exts[i] = item
			i++
		}
	}
	return exts
}

// 获取上传文件允许的最大值
func GetUploadFileSize() int64 {
	size := beego.AppConfig.DefaultString("upload_file_size", "0")
	if strings.HasSuffix(size, "MB") {
		if s, e := strconv.ParseInt(size[0:len(size)-2], 10, 64); e == nil {
			return s * 1024 * 1024
		}
	}
	if strings.HasSuffix(size, "GB") {
		if s, e := strconv.ParseInt(size[0:len(size)-2], 10, 64); e == nil {
			return s * 1024 * 1024 * 1024
		}
	}
	if strings.HasSuffix(size, "KB") {
		if s, e := strconv.ParseInt(size[0:len(size)-2], 10, 64); e == nil {
			return s * 1024
		}
	}
	if s, e := strconv.ParseInt(size, 10, 64); e == nil {
		return s * 1024
	}
	return 0
}

// 判断是否是允许上传的文件类型
func IsAllowUploadFileExt(ext string) bool {
	if strings.HasPrefix(ext, ".") {
		ext = string(ext[1:])
	}
	exts := GetUploadFileExt()
	for _, item := range exts {
		if item == "*" {
			return true
		}
		if strings.EqualFold(item, ext) {
			return true
		}
	}
	return false
}

// 下载远程文件并保存到指定位置
func DownloadAndSaveFile(remoteUrl, dstFile string) error {
	client := &http.Client{}
	uri, err := url.Parse(remoteUrl)
	if err != nil {
		return err
	}
	// Create the file
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer out.Close()

	request, err := http.NewRequest("GET", uri.String(), nil)
	request.Header.Add("Connection", "close")
	request.Header.Add("Host", uri.Host)
	request.Header.Add("Referer", uri.String())
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		_, err = io.Copy(out, resp.Body)
	} else {
		return errors.New(fmt.Sprintf("bad status: %s", resp.Status))
	}
	return nil
}

// 解压zip文件
//@param			zipFile			需要解压的zip文件
//@param			dest			需要解压到的目录
//@return			err				返回错误
func Unzip(zipFile, dest string) (err error) {
	dest = strings.TrimSuffix(dest, "/") + "/"
	// 打开一个zip格式文件
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()
	// 迭代压缩文件中的文件，打印出文件中的内容
	for _, f := range r.File {
		if !f.FileInfo().IsDir() { //非目录，且不包含__MACOSX
			if folder := dest + filepath.Dir(f.Name); !strings.Contains(folder, "__MACOSX") {
				os.MkdirAll(folder, 0777)
				if fcreate, err := os.Create(dest + strings.TrimPrefix(f.Name, "./")); err == nil {
					if rc, err := f.Open(); err == nil {
						io.Copy(fcreate, rc)
						rc.Close() //不要用defer来关闭，如果文件太多的话，会报too many open files 的错误
						fcreate.Close()
					} else {
						fcreate.Close()
						return err
					}
				} else {
					return err
				}
			}
		}
	}
	return nil
}

// 压缩文件
func Zip(source, target string) error {
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()
	source = strings.Replace(source, "\\", "/", -1)

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		path = strings.Replace(path, "\\", "/", -1)

		if path == source {
			return nil
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(strings.TrimPrefix(strings.Replace(path, "\\", "/", -1), source), "/")

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

func Compress(dst string, src string) (err error) {
	d, _ := os.Create(dst)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()

	src = strings.Replace(src, "\\", "/", -1)
	f, err := os.Open(src)

	if err != nil {
		return err
	}

	//prefix := src[strings.LastIndex(src,"/"):]
	err = compress(f, "", w)
	if err != nil {
		return err
	}

	return nil
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		if prefix != "" {
			prefix = prefix + "/" + info.Name()
		} else {
			prefix = info.Name()
		}
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		if prefix != "" {
			header.Name = prefix + "/" + header.Name
		}

		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
