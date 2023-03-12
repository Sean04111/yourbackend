package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	DB struct {
		Mysql string
	}
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	Redis struct {
		Addr string
		DB   int
	}
	ES struct {
		Addr string
	}
	Mongo struct {
		Addr string
	}
	Url struct {
		Url string
	}
	Tolerance struct {
		DBTime int64
	}
	ChatGPT struct {
		Key          string
		TargetServer string
	}
}
