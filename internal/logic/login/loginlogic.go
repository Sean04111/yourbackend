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
	gotuser, err := l.svcCtx.MysqlModel.FindOne(l.ctx, req.Email)
	switch err {
	case model.ErrNotFound:
		return &types.Loginresp{
			Status: 2,
		}, err
	case nil:
		if gotuser.Password == req.Password {
			jwttoken, err := l.GetJWT(l.svcCtx.Config.Auth.Secretkey, time.Now().Unix(), l.svcCtx.Config.Auth.Expiretime, gotuser.Uid)
			if err != nil {
				return &types.Loginresp{
					Status: 1,
				}, err
			}
			return &types.Loginresp{
				Status:      0,
				Name:        gotuser.Name,
				AccessToken: jwttoken,
				Expires:     strconv.Itoa(int(l.svcCtx.Config.Auth.Expiretime + time.Now().Unix())),
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
func (l *LoginLogic) GetJWT(key string, starttime, lasttime, uid int64) (string, error) {
	claim := make(jwt.MapClaims)
	claim["starttime"] = starttime
	claim["expiretime"] = starttime + lasttime
	claim["uid"] = uid
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claim
	return token.SignedString([]byte(key))
}
