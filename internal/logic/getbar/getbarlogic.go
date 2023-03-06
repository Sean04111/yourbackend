package getbar

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetbarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetbarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetbarLogic {
	return &GetbarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetbarLogic) Getbar(req *types.Getbarreq) (resp *types.Getbarresp, err error) {
	got, e := l.svcCtx.ArticleMysqlModel.FindOne(l.ctx, req.Id)
	if req.User == "" {

		if e != nil {
			return &types.Getbarresp{
				Status: 1,
			}, nil
		}
		return &types.Getbarresp{
			Status:  0,
			Options: []types.Eachbar{{Icon: "like.svg", Counts: got.Likes, Active: false, ActiveIcon: "likeit.svg", DisActiveIcon: "like.svg"}},
		}, nil
	} else {
		mongoclient, e := mongo.Connect(context.Background(), options.Client().ApplyURI(l.svcCtx.Config.Mongo.Addr))
		if e != nil {
			return &types.Getbarresp{
				Status: 1,
			}, nil
		}
		collection := mongoclient.Database("DB").Collection("userarticle")
		res := collection.FindOne(context.Background(), bson.M{"uid": req.User})
		var gotuser bson.M
		res.Decode(&got)
		likes := gotuser["likes"].(bson.A)
		for _, k := range likes {
			if k.(string) == req.Id {
				return &types.Getbarresp{
					Status:  0,
					Options: []types.Eachbar{{Icon: "like.svg", Counts: got.Likes, Active: true, ActiveIcon: "likeit.svg", DisActiveIcon: "like.svg"}},
				}, nil
			}
		}
		return &types.Getbarresp{
			Status:  1,
			Options: []types.Eachbar{{Icon: "like.svg", Counts: got.Likes, Active: true, ActiveIcon: "likeit.svg", DisActiveIcon: "like.svg"}},
		}, nil
	}
}
