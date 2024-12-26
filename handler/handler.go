package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//返回上传html页面

		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		// 接收文件流及存储到本地目录
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data,err %s\n", err.Error())
		}
		defer file.Close()

		fileMeta := meta
	}
}
