package config

type FileSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"Port"`
}

type ServerConfig struct {
	Name        string        `mapstructure:"name"`
	FileSrvInfo FileSrvConfig `mapstructure:"file_srv"`
}
