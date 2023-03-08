package infotable

import (
	"context"

	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfotableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfotableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfotableLogic {
	return &InfotableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfotableLogic) Infotable() (resp *types.Infotableresp, err error) {
	return &types.Infotableresp{
		Status: 0,
		Tabs: []types.SingleTab{
			{
			 TabName: "个人资料",
            ComponentName: "basicSetting",
            TabUrl: "myinfo/setting/setUserInfo",
            TabIcon: "/infoset.svg",
		},
		{
			TabName: "账号设置",
            ComponentName: "accountSetting",
            TabUrl: "myinfo/setting/setAccountInfo",
            TabIcon: "/accountset.svg",
		},
		{
			TabName: "待开发",
            ComponentName: "",
            TabUrl: "myinfo/setting",
		},
	},
	},nil
}
