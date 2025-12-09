package config

import "time"

//数据库配置

type DBConfig struct {
	Driver   string
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	Charset  string
}

//JWT 配置
type JWTConfig struct {
	SecretKey      string
	ExpirationTime time.Duration
}

//加载配置（实际项目可从环境变量或配置文件读取）
func LoadDBConfig() DBConfig {
	return DBConfig{
		Driver:   "mysql",
		Host:     "localhost",
		Port:     "3306",
		Username: "root",     //MYSQL用户名
		Password: "a5627951", //MYSQL密码
		DBName:   "blog_db",  //数据库名
		Charset:  "utf8mb4",
	}
}

func LoadJWTConfig() JWTConfig {
	return JWTConfig{
		SecretKey:      "88888888",         // 生产环境需更换为安全密钥
		ExpirationTime: 7 * 24 * time.Hour, //7天有效期
	}
}
