package updatepwd

import (
	"context"

	"yourbackend/apiService/user/api/internal/svc"
	"yourbackend/apiService/user/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
