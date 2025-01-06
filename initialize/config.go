package initialize

import (
	"cloud_storage/global"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	debug := GetEnvInfo("CLOUD_STORAGE_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("../cloud_storage/%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("../cloud_storage/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()
	// 文件路径设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		zap.S().Fatalf(err.Error())
	}

	if err := v.Unmarshal(global.ServerConfig); err != nil {
		zap.S().Fatalf(err.Error())
	}
	zap.S().Infof("配置信息: &v", global.ServerConfig)

}
