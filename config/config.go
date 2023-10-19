package config

//mapstructure 是用来读取 yaml 文件字段名 tag
type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	DBName     string `mapstructure:"dbname" json:"dbname"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

// 對應的yaml配置文件
type ServiceConfig struct {
	Port  int         `mapstructure:"port" json:"port"`
	MysqlDB    MysqlConfig `mapstructure:"mysql" json:"mysql"`
	Redis RedisConfig `mapstructure:"redis" json:"redis"`
}
