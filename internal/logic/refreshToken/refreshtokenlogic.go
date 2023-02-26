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
	Jwttoken,err:=l.GetJWT(l.svcCtx.Config.Auth.Secretkey,time.Now().Unix(),l.svcCtx.Config.Auth.Expiretime)
	if err!=nil{
		return &types.RefreshTokenresp{
			Status: 1,
		},err
	}
	return &types.RefreshTokenresp{
		Status: 0,
		Token: Jwttoken,
		Expires: strconv.Itoa(int(l.svcCtx.Config.Auth.Expiretime+time.Now().Unix())),
	},nil
}
func (l *RefreshTokenLogic)GetJWT(key string,starttime ,lasttime int64)(string ,error){
	claim:=make(jwt.MapClaims)
	claim["starttime"]=starttime
	claim["expiretime"]=starttime+lasttime
	token:=jwt.New(jwt.SigningMethodHS256)
	token.Claims=claim
	return token.SignedString([]byte(key))
}
