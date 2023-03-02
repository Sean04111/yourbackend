package searcharticle

import (
	"context"

	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/olivere/elastic/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchbytitleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchbytitleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchbytitleLogic {
	return &SearchbytitleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchbytitleLogic) Searchbytitle(req *types.Searchreq) (resp *types.Searchresp, err error) {

	return
}
func (l *SearchbytitleLogic)FromES(keyword string)(string,error){
	esclient,e:=elastic.NewClient(elastic.SetURL(l.svcCtx.Config.ES.Addr))
	if e!=nil{
		return "",e
	}
	
}
