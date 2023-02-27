package pubkey

import (
	"context"

	"yourbackend/internal/svc"
	"yourbackend/internal/types"

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
	return &types.Pubkeyresp{
		Status: 0,
		Pubkey: string(l.svcCtx.RsaOps.GetPubkey()),
	},nil
}
