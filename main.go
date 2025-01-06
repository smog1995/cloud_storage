package main

import (
	"cloud_storage/global"
	"cloud_storage/initialize"
	"fmt"
	"go.uber.org/zap"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	router := initialize.Routers()
	router.LoadHTMLGlob("static/view/*")
	zap.S().Debugf("启动服务器, 端口： %d", global.ServerConfig.Port)
	if err := router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}
}
