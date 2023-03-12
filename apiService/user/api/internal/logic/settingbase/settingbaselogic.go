package settingbase

import (
	"context"

	"yourbackend/apiService/user/api/internal/svc"
	"yourbackend/apiService/user/api/internal/types"

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

func (l *SettingbaseLogic) Settingbase(req *types.Settingbasereq) (resp *types.Settingbaseresp, err error) {
	// todo: add your logic here and delete this line

	return
}
