package svc

import (
	"yourbackend/internal/config"
	"yourbackend/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	MysqlModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		MysqlModel: model.NewUserModel(sqlx.NewMysql(c.DB.Mysql)),
	}
}
