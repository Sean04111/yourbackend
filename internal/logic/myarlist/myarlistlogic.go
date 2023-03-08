package myarlist

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"yourbackend/internal/model"

	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MyarlistLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	Catcher   []string
	MyCatcher []string
}

func NewMyarlistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MyarlistLogic {
	return &MyarlistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MyarlistLogic) Myarlist(req *types.Myarlistreq) (resp *types.Myarlistresp, err error) {
	e := l.FromMongo()
	if e != nil {
		return &types.Myarlistresp{
			Status: 1,
		}, nil
	}
	if req.Input == "" {
		var res []types.Arti
		for _, k := range l.MyCatcher {
			r, _ := l.svcCtx.ArticleMysqlModel.FindOne(l.ctx, k)
			res = append(res, types.Arti{
				Id:          r.Mongoid,
				Title:       r.Title,
				Description: r.Fewcontent,
				Likes:       int(r.Likes),
				Reads:       int(r.Views),
				Url:         r.Url,
				Pubtime:     r.Pubtime,
				Imglink:     r.Coverlinks,
			})
		}
		return &types.Myarlistresp{
			Status: 0,
			List:   res,
		}, nil
	} else {
		res, er := l.FromES(req.Input)
		if er != nil {

			return &types.Myarlistresp{
				Status: 2,
			}, nil
		}
		if len(l.Catcher) > 20 {
			return &types.Myarlistresp{
				Status: 3,
			}, nil
		}
		var modelres []types.Arti
		for i := 0; i < len(res); i++ {
			for _, k := range l.MyCatcher {
				if k == res[i].Mongoid {
					modelres = append(modelres, types.Arti{
						Id:          res[i].Mongoid,
						Title:       res[i].Title,
						Description: res[i].Fewcontent,
						Likes:       int(res[i].Likes),
						Reads:       int(res[i].Views),
						Url:         res[i].Url,
						Pubtime:     res[i].Pubtime,
						Imglink:     res[i].Coverlinks,
					})
				}
			}
		}
		return &types.Myarlistresp{
			Status: 0,
			List:   modelres,
		}, nil
	}
}

type Artic struct {
	Fewcontent string
	Mongoid    string
	Title      string
}

func (l *MyarlistLogic) FromES(keyword string) ([]*model.Articles, error) {
	esclient, e := elastic.NewClient(elastic.SetURL(l.svcCtx.Config.ES.Addr), elastic.SetSniff(false))
	if e != nil {
		return nil, e
	}
	matchquery := elastic.NewMatchQuery("fewcontent", keyword)
	searchresult, er := esclient.Search().Index("article").Query(matchquery).From(0).Size(10).Do(context.TODO())
	if er != nil {
		return nil, er
	} else if searchresult.TotalHits() > 0 {
		var got Artic
		for _, item := range searchresult.Each(reflect.TypeOf(got)) {
			if t, ok := item.(Artic); ok {
				l.Catcher = append(l.Catcher, t.Mongoid)
			}
		}
		var Res []*model.Articles
		for _,k:=range l.Catcher{
			res, e := l.svcCtx.ArticleMysqlModel.FindOne(l.ctx, k)
			if e != nil {
				fmt.Println("这里错了：",e)
				return nil, e
			}
			Res = append(Res, res)
		}
		return Res, nil
	} else {
		return nil, errors.New("NoMatch")
	}
}
func (l *MyarlistLogic) FromMongo() error {
	mongoclient, e := mongo.Connect(context.Background(), options.Client().ApplyURI(l.svcCtx.Config.Mongo.Addr))
	if e != nil {
		return e
	}
	collection := mongoclient.Database("DB").Collection("userarticle")
	res := collection.FindOne(context.Background(), bson.M{"uid": l.ctx.Value("uid")})
	var got bson.M
	res.Decode(&got)
	artis := got["articles"].(bson.A)
	for _, k := range artis {
		l.MyCatcher = append(l.MyCatcher, k.(string))
	}
	return nil
}
