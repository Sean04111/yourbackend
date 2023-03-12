package getdraft

import (
	"context"

	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GetdraftLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetdraftLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetdraftLogic {
	return &GetdraftLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetdraftLogic) Getdraft(req *types.Getdraftreq) (resp *types.Getdraftresp, err error) {
	mongoclient,e:=mongo.Connect(context.Background(),options.Client().ApplyURI(l.svcCtx.Config.Mongo.Addr))
	if e!=nil{
		return &types.Getdraftresp{
			Status: 1,
		},nil
	}
	collection:=mongoclient.Database("DB").Collection("article")
	res:=collection.FindOne(context.Background(),bson.M{"arid":req.Arid})
	var got bson.M
	res.Decode(&got)
	if got["authorid"].(string)!=l.ctx.Value("uid"){
		return &types.Getdraftresp{
			Status: 2,
		},nil
	}
	return &types.Getdraftresp{
		Status: 0,
		Text:got["content"].(string),
		Id:got["arid"].(string),
		IsPublish: got["ispublish"].(bool),
	},nil
}
