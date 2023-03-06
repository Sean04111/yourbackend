package tools

import (
	"context"

	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ToolsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewToolsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToolsLogic {
	return &ToolsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *ToolsLogic) Tools() (resp *types.Toolsresp, err error) {
	return &types.Toolsresp{
		Status: 0,
		Items: []types.Tool{{Url: "https://tbghg.top/",Title: "学长博客"},},
	}, nil
}
