package articleread

import (
	"context"
	"strconv"
	"sync"
	"time"

	"yourbackend/internal/model"
	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArticlereadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticlereadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticlereadLogic {
	return &ArticlereadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticlereadLogic) Articleread(req *types.Readarticlereq) (resp *types.Readarticleresp, err error) {
	arbm, er := l.FromArtiMongo(req.Arid)
	if er != nil {
		return &types.Readarticleresp{
			Status: 1,
		}, nil
	}
	if req.Userid == "" {
		wg := new(sync.WaitGroup)
		wg.Add(2)
		errchan := make(chan error, 2)
		go func() {
			e := l.Artihandler(arbm)
			if e != nil {
				errchan <- e
			}
			wg.Done()
		}()
		go func() {
			errr := l.Userhanlder(arbm)
			if errr != nil {
				errchan <- errr
			}
			wg.Done()
		}()
		wg.Wait()
		select {
		case <-errchan:

			return &types.Readarticleresp{
				Status: 1,
			}, nil
		default:
			return &types.Readarticleresp{
				Status:  0,
				Content: arbm["content"].(string),
				Title:   arbm["title"].(string),
				IsEdit:  false,
			}, nil
		}
	} else {
		var belong = false
		wg := new(sync.WaitGroup)
		wg.Add(3)
		errchan := make(chan error, 3)
		go func() {
			userbm, e := l.FromUserMongo(req.Userid)
			artilist := userbm["articles"].(bson.A)
			for _, k := range artilist {
				if k == req.Arid {
					belong = true
				}
			}
			errchan <- e
			wg.Done()
		}()
		go func() {
			errr := l.Artihandler(arbm)
			errchan <- errr
			wg.Done()
		}()
		go func() {
			e := l.Userhanlder(arbm)
			errchan <- e
			wg.Done()
		}()
		wg.Wait()
		select {
		case <-errchan:
			return &types.Readarticleresp{
				Status: 1,
			}, nil
		default:
			return &types.Readarticleresp{
				Status:  0,
				Content: arbm["content"].(string),
				Title:   arbm["title"].(string),
				IsEdit:  belong,
			}, nil
		}
	}
}
func (l *ArticlereadLogic) FromUserMongo(uid string) (bson.M, error) {
	mongoclient, e := mongo.Connect(context.Background(), options.Client().ApplyURI(l.svcCtx.Config.Mongo.Addr))
	if e != nil {
		return nil, e
	}
	collection := mongoclient.Database("DB").Collection("userarticle")
	res := collection.FindOne(context.Background(), bson.M{"uid": uid})
	var ans bson.M
	if e := res.Decode(&ans); e != nil {
		return nil, e
	} else {
		return ans, nil
	}
}
func (l *ArticlereadLogic) FromArtiMongo(arid string) (bson.M, error) {
	mongoclient, e := mongo.Connect(context.Background(), options.Client().ApplyURI(l.svcCtx.Config.Mongo.Addr))
	if e != nil {
		return nil, e
	}
	collection := mongoclient.Database("DB").Collection("article")
	res := collection.FindOne(context.Background(), bson.M{"arid": arid})
	var ans bson.M
	if e := res.Decode(&ans); e != nil {
		return nil, e
	} else {
		return ans, nil
	}
}
func (l *ArticlereadLogic) Artihandler(article bson.M) error {
	mongoclient, e := mongo.Connect(context.Background(), options.Client().ApplyURI(l.svcCtx.Config.Mongo.Addr))
	if e != nil {
		return e
	}
	articlecollection := mongoclient.Database("DB").Collection("article")
	daysdata := article["daysdata"].(bson.A)      //
	lastrefresh := article["lastrefresh"].(int64) //
	now := time.Now().Unix()
	if (now/24/3600 - lastrefresh/24/3600) > 0 {
		for i := 1; i < 7; i++ {
			daysdata[i-1] = daysdata[i]
		}
		daysdata[6] = 1
	} else {
		daysdata[6] = daysdata[6].(int64) + 1
	}
	_, er := articlecollection.UpdateOne(context.Background(), bson.M{"arid": article["arid"].(string)}, bson.M{"$set": bson.M{"daysdata": daysdata, "views": article["views"].(int64) + 1, "lastrefresh": time.Now().Unix()}})
	if er != nil {
		return er
	}
	err := l.svcCtx.ArticleMysqlModel.Update(l.ctx, &model.Articles{
		Mongoid:    article["arid"].(string),
		Title:      article["title"].(string),
		Fewcontent: article["fewcontent"].(string),
		Likes:      article["likes"].(int64),
		Views:      article["views"].(int64) + 1,
		Url:        article["url"].(string),
		Pubtime:    strconv.Itoa(int(article["created"].(int64))),
		Coverlinks: article["coverlink"].(string),
	})
	if err != nil {
		return err
	}
	return nil
}
func (l *ArticlereadLogic) Userhanlder(article bson.M) error {
	mongoclient, e := mongo.Connect(context.Background(), options.Client().ApplyURI(l.svcCtx.Config.Mongo.Addr))
	if e != nil {
		return e
	}
	usercollection := mongoclient.Database("DB").Collection("userarticle")
	articlecollection := mongoclient.Database("DB").Collection("article")
	res := usercollection.FindOne(context.Background(), bson.M{"uid": article["authorid"].(string)})
	var gotuser bson.M
	res.Decode(&gotuser)
	lastrefresh := gotuser["lastrefresh"].(int64)
	now := time.Now().Unix()
	alldata := gotuser["alldata"].(bson.A)
	if (now/24/3600 - lastrefresh/24/3600) > 0 {
		for i := 1; i < 7; i++ {
			alldata[i-1] = alldata[i]
		}
		alldata[6] = 1

	} else {
		allarticles := gotuser["articles"].(bson.A)
		var alltoday int64
		for _, k := range allarticles {
			res := articlecollection.FindOne(context.Background(), bson.M{"arid": k.(string)})
			var ans bson.M
			res.Decode(&ans)
			alltoday = alltoday + ans["daysdata"].(bson.A)[6].(int64)
		}
		alldata[6] = alltoday
	}
	_, er := usercollection.UpdateOne(context.Background(), bson.M{"uid": gotuser["uid"]}, bson.M{"$set": bson.M{"alldata": alldata, "lastrefresh": time.Now().Unix()}})
	if er != nil {
		return er
	} else {
		return nil
	}
}
