package searcharticle

import (
	"context"
	"errors"
	"reflect"

	"yourbackend/internal/model"
	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/olivere/elastic/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchbytitleLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	Catcher []string
}

func NewSearchbytitleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchbytitleLogic {
	return &SearchbytitleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchbytitleLogic) Searchbytitle(req *types.Searchreq) (resp *types.Searchresp, err error) {
	res, er := l.FromES(req.Articlename)
	if er != nil {
		return &types.Searchresp{
			Status: 1,
		}, nil
	}
	if len(l.Catcher)>20{
		return &types.Searchresp{
			Status: 1,
		},nil
	}
	var modelres []types.Article
	for i := 0; i < len(res); i++ {
		modelres = append(modelres, types.Article{
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
	return &types.Searchresp{
		Status:      0,
		Articlelist: modelres,
	}, nil
}

type Arti struct {
	Fewcontent string
	Mongoid    string
	Title      string
}

func (l *SearchbytitleLogic) FromES(keyword string) ([]*model.Articles, error) {
	esclient, e := elastic.NewClient(elastic.SetURL(l.svcCtx.Config.ES.Addr), elastic.SetSniff(false))
	if e != nil {
		return nil, e
	}
	matchquery := elastic.NewMatchQuery("fewcontent", keyword)
	searchresult, er := esclient.Search().Index("article").Query(matchquery).Do(context.Background())
	if er != nil {
		return nil, er
	}else 
	if searchresult.TotalHits() > 0 {
		var got Arti
		for _, item := range searchresult.Each(reflect.TypeOf(got)) {
			if t, ok := item.(Arti); ok {
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
			return Res, nil
		}
	} else {
		return nil, errors.New("NoMatch")
	}
	return nil,nil
}

