package config

const (
	// MySQLSource : 要连接的数据库源；
	// 其中test:test 是用户名密码；
	// 127.0.0.1:3306 是ip及端口；
	// fileserver 是数据库名;
	// charset=utf8 指定了数据以utf8字符编码进行传输
	MySQLSource = "root:root1234@tcp(127.0.0.1:3306)/cloud_storage?charset=utf8&parseTime=True"
)
