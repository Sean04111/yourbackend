package updatepwd

import (
	"context"

	"yourbackend/internal/model"
	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logx"
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
	if er != nil {
		return &types.Updatepwdresp{
			Status: 1,
		}, er
	}
	if realcode == req.Code {
		newuser, e := l.svcCtx.MysqlModel.FindOne(l.ctx, req.Email)
		if e != nil {
			return &types.Updatepwdresp{
				Status: 1,
			}, e
		}
		newuser.Password = req.Password
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
			Status: 2,
		}, nil
	}
}
func (l *UpdatepwdLogic) FromRedis(email string) (string, error) {
	redisclient := redis.NewClient(&redis.Options{
		Addr: "121.36.131.50:6379",
		DB:   0,
	})
	return redisclient.Get(email).Result()
}
