package svc

import (
	rsarand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"yourbackend/internal/config"
	"yourbackend/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type rsaoperator struct {

}
func (r *rsaoperator) Decode(cipher []byte) []byte {
	f, e := ioutil.ReadFile("internal/svc/private.pem")
	if e != nil {
		panic(e)
	}
	block, _ := pem.Decode(f)
	privatekey, er := x509.ParsePKCS1PrivateKey(block.Bytes)
	if er != nil {
		panic(er)
	}
	real, err := rsa.DecryptPKCS1v15(rsarand.Reader, privatekey, cipher)
	if err != nil {
		panic(err)
	}
	return real
}
func(r *rsaoperator)GetPubkey()[]byte{
	f, e := ioutil.ReadFile("internal/svc/public.pem")
	if e != nil {
		panic(e)
	}
	return f
}
//to dynamically generate keys
func(r *rsaoperator)generatekey(){}

type ServiceContext struct {
	Config     config.Config
	MysqlModel model.UserModel
	RsaOps     rsaoperator
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		MysqlModel: model.NewUserModel(sqlx.NewMysql(c.DB.Mysql)),
	}
}
