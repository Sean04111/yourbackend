package refreshToken

import (
	"context"
	"strconv"
	"time"

	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenreq) (resp *types.RefreshTokenresp, err error) {
	gotuser,err:=l.svcCtx.MysqlModel.FindOne(l.ctx,req.Email)
	if err!=nil{
		return &types.RefreshTokenresp{
			Status: 1,
		},err
	}
	Jwttoken,e:=l.GetJWT(l.svcCtx.Config.Auth.AccessSecret,req.Email,gotuser.Uid,time.Now().Unix(),l.svcCtx.Config.Auth.AccessExpire)
	if e!=nil{
		return &types.RefreshTokenresp{
			Status: 1,
		},e
	}
	return &types.RefreshTokenresp{
		Status: 0,
		Name: gotuser.Name,
		Token: Jwttoken,
		Expires: strconv.Itoa(int(l.svcCtx.Config.Auth.AccessExpire+time.Now().Unix())),
	},nil
}
func (l *RefreshTokenLogic) GetJWT(key,email,uid string, starttime, lasttime int64) (string, error) {
	claim := make(jwt.MapClaims)
	claim["starttime"] = starttime
	claim["expiretime"] = starttime + lasttime
	claim["email"] = email
	claim["uid"]=uid
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claim
	return token.SignedString([]byte(key))
}

