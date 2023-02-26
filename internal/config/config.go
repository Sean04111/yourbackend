package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	DB struct{
		Mysql string
	}
	Auth struct{
		Secretkey string
		Expiretime int64
	}
}
