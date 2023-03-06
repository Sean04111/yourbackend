package myarlist

import (
	"context"
	"errors"
	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"yourbackend/internal/model"

	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MyarlistLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	Catcher []string
}

func NewMyarlistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MyarlistLogic {
	return &MyarlistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MyarlistLogic) Myarlist(req *types.Myarlistreq) (resp *types.Myarlistresp, err error) {
	res, er := l.FromES(req.Input)
	if er != nil {
		return &types.Myarlistresp{
			Status: 1,
		}, nil
	}
	if len(l.Catcher) > 20 {
		return &types.Myarlistresp{
			Status: 1,
		}, nil
	}
	var modelres []types.Arti
	mongoclient, e := mongo.Connect(context.Background(), options.Client().ApplyURI(l.svcCtx.Config.Mongo.Addr))
	if e != nil {
		return &types.Myarlistresp{
			Status: 1,
		}, nil
	}
	r := mongoclient.Database("DB").Collection("article")
	for i := 0; i < len(res); i++ {
		re := r.FindOne(context.Background(), bson.M{"arid": res[i].Mongoid})
		var got bson.M
		re.Decode(&got)
		if got["authorid"].(string) == l.ctx.Value("uid") {
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
	return &types.Myarlistresp{
		Status: 0,
		List:   modelres,
	}, nil
	
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
		for i := 0; i < len(l.Catcher); i++ {
			res, e := l.svcCtx.ArticleMysqlModel.FindOne(l.ctx, l.Catcher[i])
			if e != nil {
				return nil, e
			}
			Res = append(Res, res)
		}
		return Res, nil
	} else {
		return nil, errors.New("NoMatch")
	}
}
