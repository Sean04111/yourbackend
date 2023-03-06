package getsingledata

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetsingledataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetsingledataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetsingledataLogic {
	return &GetsingledataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetsingledataLogic) Getsingledata(req *types.Getsingledatareq) (resp *types.Getsingledataresp, err error) {
	mongoclient, e := mongo.Connect(context.Background(), options.Client().ApplyURI(l.svcCtx.Config.Mongo.Addr))
	if e != nil {
		return &types.Getsingledataresp{
			Status: 1,
		}, nil
	}
	collection := mongoclient.Database("DB").Collection("article")
	res := collection.FindOne(context.Background(), bson.M{"arid": req.Id})
	var got bson.M
	er := res.Decode(&got)
	if er != nil {
		return &types.Getsingledataresp{
			Status: 1,
		}, nil
	}
	daysdata := got["daysdata"].(bson.A)
	var catcher []int64
	for _, k := range daysdata {
		catcher = append(catcher, k.(int64))
	}
	return &types.Getsingledataresp{
		Status:    1,
		LineData:  catcher,
		LineLable: []string{"6天前", "5天前", "4天前", "3天前", "2天前", "昨天", "今天"},
	}, nil
}
