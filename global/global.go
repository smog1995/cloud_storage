package global

import (
	"cloud_storage/config"
)

// 全局可使用的变量
var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
)
