package main

import (
	"cloud_storage/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("static/view/*")
	router.GET("/file/upload", handler.UploadHandler)
	router.POST("/file/upload", handler.UploadHandler)
	router.Run(":8080")
	//http.HandleFunc("/file/upload", handler.UploadHandler)
	//err := http.ListenAndServe(":8080", nil)
	//if err != nil {
	//	fmt.Printf("Fail to start server, err%s", err.Error())
	//}
}
