package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 处理文件上传
func UploadHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "文件上传",
		})

	} else if c.Request.Method == "POST" {
		file, header, err := c.Request.FormFile("file")
		defer file.Close()
		if err != nil {
			fmt.Printf("Failed to get data,err %s\n", err.Error())
			return
		}
		fmt.Printf(header.Filename)

	}

}
