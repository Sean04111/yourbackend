package settingbase

import (
	"context"

	"yourbackend/internal/model"
	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SettingbaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSettingbaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SettingbaseLogic {
	return &SettingbaseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SettingbaseLogic) Settingbase(req *types.Settingbasereq) (*types.Settingbaseresp, error) {
	gotuser, err := l.svcCtx.MysqlModel.FindOne(l.ctx, l.ctx.Value("email").(string))
	if err != nil {
		return &types.Settingbaseresp{
			Status: 1,
		}, nil
	}
	gotuser.Name = req.Name
	gotuser.Profession = req.Profession
	gotuser.Type = req.Type
	 e := l.svcCtx.MysqlModel.Update(l.ctx, &model.User{
		Uid:gotuser.Uid,
		Email: gotuser.Email,
		Password: gotuser.Password,
		Name:gotuser.Name,
		AvatarLink: gotuser.AvatarLink,
		Profession: gotuser.Profession,
		Type: gotuser.Type,
	})
	if e != nil {
		return &types.Settingbaseresp{
			Status: 1,
		}, nil
	} else {
		return &types.Settingbaseresp{
			Status: 0,
		},nil

	}

}
