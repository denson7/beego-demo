package controllers

import (
	"beego-demo/utils"
	"beego-demo/utils/cryptil"
	"beego-demo/utils/filetil"
	"github.com/astaxie/beego"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type UploadController struct {
	BaseController
}

type Attachment struct {
	FileName   string    `json:"file_name"`
	FilePath   string    `json:"file_path"`
	HttpPath   string    `json:"http_path"`
	FileSize   float64   `json:"file_size"`
	FileExt    string    `json:"file_ext"`
}

// 上传附件或图片
func (c *UploadController) Upload() {
	name := "file"
	file, moreFile, err := c.GetFile(name)
	if err == http.ErrMissingFile || moreFile == nil {
		c.JsonResult(6003, "没有发现需要上传的文件", "")
	}

	if err != nil {
		c.JsonResult(6002, err.Error(), "")
	}

	defer file.Close()

	type Size interface {
		Size() int64
	}

	if utils.GetUploadFileSize() > 0 && moreFile.Size > utils.GetUploadFileSize() {
		c.JsonResult(6009, "文件大小超过了限定的最大值", "")
	}

	ext := filepath.Ext(moreFile.Filename)
	// 文件必须带有后缀名
	if ext == "" {
		c.JsonResult(6003, "无法解析文件的格式", "")
	}
	// 如果文件类型设置为 * 标识不限制文件类型
	if utils.IsAllowUploadFileExt(ext) == false {
		c.JsonResult(6004, "不允许的文件类型", "")
	}

	fileName := cryptil.UniqueId()
	filePath := filepath.Join("./", "uploads")

	// 将图片和文件分开存放
	if filetil.IsImageExt(moreFile.Filename) {
		filePath = filepath.Join(filePath, "images", fileName+ext)
	} else {
		filePath = filepath.Join(filePath, "files", fileName+ext)
	}

	path := filepath.Dir(filePath)
	_ = os.MkdirAll(path, os.ModePerm)
	err = c.SaveToFile(name, filePath)

	if err != nil {
		beego.Error("保存文件失败 -> ", err)
		c.JsonResult(6005, "保存文件失败", "")
	}

	attachment := &Attachment{}
	attachment.FileName = moreFile.Filename
	attachment.FileExt = ext
	attachment.FilePath = strings.TrimPrefix(filePath, "./")

	if fileInfo, err := os.Stat(filePath); err == nil {
		attachment.FileSize = float64(fileInfo.Size())
	}
	if filetil.IsImageExt(moreFile.Filename) {
		attachment.HttpPath = "/" + strings.Replace(strings.TrimPrefix(filePath, "./"), "\\", "/", -1)
		if strings.HasPrefix(attachment.HttpPath, "//") {
			attachment.HttpPath = string(attachment.HttpPath[1:])
		}
	}

	result := map[string]interface{}{
		"attach": attachment,
	}

	c.JsonResult(200, "ok", result)
}
