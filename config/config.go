package config

type ServerConfig struct {
	Name           string `mapstructure:"name"`
	FileLocation   string `mapstructure:"file_location"`
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	CephAccessKey  string `mapstructure:"ceph_access_key"`
	CephSecretKey  string `mapstructure:"ceph_secret_key"`
	CephGWEndpoint string `mapstructure:"ceph_gateway_endpoint"`
}
