package getardata

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetardataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetardataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetardataLogic {
	return &GetardataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetardataLogic) Getardata() (resp *types.Getardataresp, err error) {
	uid := l.ctx.Value("uid").(string)
	mongoclient, e := mongo.Connect(context.Background(), options.Client().ApplyURI(l.svcCtx.Config.Mongo.Addr))
	if e != nil {
		return &types.Getardataresp{
			Status: 1,
		}, nil
	}
	collection := mongoclient.Database("DB").Collection("userarticle")
	res := collection.FindOne(context.TODO(), bson.M{"uid": uid})
	var got bson.M
	er := res.Decode(&got)
	if er != nil {
		return &types.Getardataresp{
			Status: 1,
		}, nil
	}
	datalist := got["alldata"].(bson.A)
	var datas []int64
	for _, k := range datalist {
		datas = append(datas, int64(k.(int32)))
	}
	return &types.Getardataresp{
		Status:   0,
		Bardata:  datas,
		Barlabel: []string{"6天前", "5天前", "4天前", "3天前", "2天前", "昨天", "今天"},
	}, nil
}
