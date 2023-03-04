package login

import (
	"context"
	"strconv"
	"time"

	"yourbackend/internal/model"
	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.Loginreq) (resp *types.Loginresp, err error) {
	requsetpassword,e:=l.svcCtx.RsaOps.Decode(req.Password)
	if e!=nil{
		return &types.Loginresp{
			Status: 1,
		},nil
	}
	gotuser, err := l.svcCtx.MysqlModel.FindOne(l.ctx, req.Email)
	switch err {
	case model.ErrNotFound:
		return &types.Loginresp{
			Status: 2,
		}, err
	case nil:
		if bcrypt.CompareHashAndPassword([]byte(gotuser.Password),requsetpassword)==nil {
			jwttoken, err := l.GetJWT(l.svcCtx.Config.Auth.AccessSecret,req.Email,gotuser.Uid, time.Now().Unix(), l.svcCtx.Config.Auth.AccessExpire)
			if err != nil {
				return &types.Loginresp{
					Status: 1,
				}, err
			}
			return &types.Loginresp{
				Status:      0,
				Name:        gotuser.Name,
				AccessToken: jwttoken,
				Expires:     strconv.Itoa(int(l.svcCtx.Config.Auth.AccessExpire + time.Now().Unix())),
			}, nil
		} else {
			return &types.Loginresp{
				Status: 3,
			}, nil
		}
	default:
		return &types.Loginresp{
			Status: 1,
		}, err
	}
}
func (l *LoginLogic) GetJWT(key,email,uid string, starttime, lasttime  int64) (string, error) {
	claim := make(jwt.MapClaims)
	claim["starttime"] = starttime
	claim["expiretime"] = starttime + lasttime
	claim["email"] = email
	claim["uid"]=uid
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claim
	return token.SignedString([]byte(key))
}
