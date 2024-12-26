package global

import (
	"cloud_storage/config"
)

var (
	ServerConfig *config.ServerConfig  = &config.ServerConfig{}
	FileConfig   *config.FileSrvConfig = &config.FileSrvConfig{}
)
