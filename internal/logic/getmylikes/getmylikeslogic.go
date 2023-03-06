package getmylikes

import (
	"context"

	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GetmylikesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetmylikesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetmylikesLogic {
	return &GetmylikesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetmylikesLogic) Getmylikes() (resp *types.Getmylikesresp, err error) {
	mongoclient,e:=mongo.Connect(context.Background(),options.Client().ApplyURI(l.svcCtx.Config.Mongo.Addr))
	if e!=nil{
		return &types.Getmylikesresp{
			Status: 1,
		},nil
	}
	collection:=mongoclient.Database("DB").Collection("userarticle")
	res:=collection.FindOne(context.Background(),bson.M{"uid":l.ctx.Value("uid").(string)})
	var got bson.M
	er:=res.Decode(&got)
	if er!=nil{
		return &types.Getmylikesresp{
			Status:1,
		},nil
	}
	var ans []types.Ar
	for _,k:=range got["articles"].(bson.A){
		res,err:=l.svcCtx.ArticleMysqlModel.FindOne(l.ctx,k.(string))
		if err!=nil{
			return &types.Getmylikesresp{
				Status: 1,
			},nil
		}
		ans = append(ans, types.Ar{
			ArticaName: res.Title,
			ArticalLink: res.Url,
			ArticalImgLink: res.Coverlinks,
			ArticalID: res.Mongoid,
		})
	}
	return &types.Getmylikesresp{
		Status: 0,
		Datalikes: ans,
	},nil
}
