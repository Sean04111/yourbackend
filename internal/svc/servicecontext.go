package svc

import (
	rsarand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"yourbackend/internal/config"
	"yourbackend/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type rsaoperator struct {

}
func (r *rsaoperator) Decode(cipher string) ([]byte,error ){
	cipherbyte,_:=base64.StdEncoding.DecodeString(cipher)
	f, e := ioutil.ReadFile("internal/svc/private.pem")
	if e != nil {
		panic(e)
	}
	block, _ := pem.Decode(f)
	privatekey, er := x509.ParsePKCS1PrivateKey(block.Bytes)
	if er != nil {
		panic(er)
	}
	return rsa.DecryptPKCS1v15(rsarand.Reader, privatekey, cipherbyte)
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
	ArticleMysqlModel model.ArticlesModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		MysqlModel: model.NewUserModel(sqlx.NewMysql(c.DB.Mysql)),
		ArticleMysqlModel: model.NewArticlesModel(sqlx.NewMysql(c.DB.Mysql)),
	}
}
