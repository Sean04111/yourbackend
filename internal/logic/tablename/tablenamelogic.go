package tablename

import (
	"context"

	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TablenameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTablenameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TablenameLogic {
	return &TablenameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TablenameLogic) Tablename() (resp *types.Tablenameresp, err error) {
	return &types.Tablenameresp{
		Status:     0,
		Chartlable: "文章总数据",
	}, nil
	return
}
