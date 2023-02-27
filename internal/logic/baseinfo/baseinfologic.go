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

func (l *BaseinfoLogic) Baseinfo() (resp *types.Baseinfor, err error) {
	return &types.Baseinfor{
		Info: types.Info{
			Usermail: l.ctx.Value("email").(string),
		},
	},nil
}
