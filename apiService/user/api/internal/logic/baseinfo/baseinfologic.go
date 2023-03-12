package baseinfo

import (
	"context"

	"yourbackend/apiService/user/api/internal/svc"
	"yourbackend/apiService/user/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
