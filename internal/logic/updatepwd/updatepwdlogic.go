package updatepwd

import (
	"context"

	"yourbackend/internal/model"
	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type UpdatepwdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatepwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatepwdLogic {
	return &UpdatepwdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatepwdLogic) Updatepwd(req *types.Updatepwdreq) (resp *types.Updatepwdresp, err error) {
	realcode, er := l.FromRedis(req.Email)
	frontendcode, e := l.FromRedis(realcode)
	if er != nil || e != nil {
		return &types.Updatepwdresp{
			Status: 1,
		}, er
	}
	if realcode == req.Code {
		if req.Check == frontendcode {
			newuser, e := l.svcCtx.MysqlModel.FindOne(l.ctx, req.Email)
			if e != nil {
				return &types.Updatepwdresp{
					Status: 1,
				}, e
			}
			stringpwd,errorr:=l.svcCtx.RsaOps.Decode(req.Password)
			if errorr!=nil{
				return &types.Updatepwdresp{
					Status: 1,
				},nil
			}
			storepwd, erro := bcrypt.GenerateFromPassword(stringpwd, 10)
			if erro != nil {
				return &types.Updatepwdresp{
					Status: 1,
				},erro
			}
			newuser.Password = string(storepwd)
			err := l.svcCtx.MysqlModel.Update(l.ctx, &model.User{
				Uid:      newuser.Uid,
				Email:    newuser.Email,
				Password: newuser.Password,
				Name:     newuser.Name,
			})
			if err != nil {
				return &types.Updatepwdresp{
					Status: 1,
				}, err
			} else {
				return &types.Updatepwdresp{
					Status: 0,
				}, nil
			}
		} else {
			return &types.Updatepwdresp{
				Status: 3,
			}, nil
		}
	} else {
		return &types.Updatepwdresp{
			Status: 2,
		}, nil
	}
}
func (l *UpdatepwdLogic) FromRedis(keyword string) (string, error) {
	redisclient := redis.NewClient(&redis.Options{
		Addr: l.svcCtx.Config.Redis.Addr,
		DB:   l.svcCtx.Config.Redis.DB,
	})
	return redisclient.Get(keyword).Result()
}
