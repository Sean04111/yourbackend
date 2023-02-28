package baseinfo

import (
	"context"

	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BaseinfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBaseinfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseinfoLogic {
	return &BaseinfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BaseinfoLogic) Baseinfo() (*types.Baseinfor, error) {
	gotuser, err := l.svcCtx.MysqlModel.FindOne(l.ctx, l.ctx.Value("email").(string))
	if err != nil {
		return &types.Baseinfor{
			Status: 1,
		}, nil
	} else {
		return &types.Baseinfor{
			Status: 0,
			Info: types.Info{
				AvatarLink: gotuser.AvatarLink,
				UserName: gotuser.Name,
				Profession: gotuser.Profession,
				Type: gotuser.Type,
				Usermail: l.ctx.Value("email").(string),
			},
		}, nil
	}

}
