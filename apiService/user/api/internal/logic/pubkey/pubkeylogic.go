package pubkey

import (
	"context"

	"yourbackend/apiService/user/api/internal/svc"
	"yourbackend/apiService/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PubkeyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPubkeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PubkeyLogic {
	return &PubkeyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PubkeyLogic) Pubkey() (resp *types.Pubkeyresp, err error) {
	// todo: add your logic here and delete this line

	return
}
