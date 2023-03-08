package updatecontent

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"yourbackend/internal/model"
	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/bwmarrin/snowflake"
	"github.com/olivere/elastic/v7"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UpdatecontentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatecontentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatecontentLogic {
	return &UpdatecontentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

//frequent operation !
// improvement needed!
//redis should be here to catch the request !!

func (l *UpdatecontentLogic) Updatecontent(req *types.Articlereq) (*types.Articleresp, error) {
	gotuser, e := l.svcCtx.MysqlModel.FindOne(l.ctx, l.ctx.Value("email").(string))
	if req.Arid == "" {
		if e != nil {
			return &types.Articleresp{
				Status: 1,
			}, nil
		}
		bson, arid, title, fewcontent := l.BsonMFiller(req.Content, gotuser.AvatarLink, gotuser.Name, "", time.Now().Unix(), req.IsPublish, req.Content == "")
		wg := new(sync.WaitGroup)
		wg.Add(2)
		errchan := make(chan error, 2)
		go func() {
			defer wg.Done()
			err := l.ToES(arid, title, fewcontent, "insert", req.Content == "")
			if err != nil {
				errchan <- err
			} else {
				return
			}
		}()
		go func() {
			defer wg.Done()
			er1, er2 := l.ToMongo(bson, "insert", req)
			if er1 != nil {
				errchan <- er1
			}
			if er2 != nil {
				errchan <- er2
			} else {
				return
			}
		}()
		wg.Wait()
		select {
		case <-errchan:
			return &types.Articleresp{
				Status: 1,
			}, nil
		default:
			return &types.Articleresp{
				Status:    0,
				Arid:      arid,
				IsPublish: req.IsPublish,
			}, nil
		}
	} else {
		bson, arid, title, fewcontent := l.BsonMFiller(req.Content, gotuser.AvatarLink, gotuser.Name, req.Arid, time.Now().Unix(), req.IsPublish, req.Content == "")
		errchan := make(chan error, 3)
		wg := new(sync.WaitGroup)
		wg.Add(2)
		go func() {
			defer wg.Done()
			err1, err2 := l.ToMongo(bson, "update", req)
			if err1 != nil {
				errchan <- err1
			}
			if err2 != nil {
				errchan <- err2
			} else {
				return
			}
		}()
		go func() {
			defer wg.Done()
			err := l.ToES(arid, title, fewcontent, "update", req.Content == "")
			if err != nil {
				errchan <- err
			} else {
				return
			}
		}()
		wg.Wait()
		select {
		case <-errchan:
			return &types.Articleresp{
				Status: 1,
			}, nil
		default:
			return &types.Articleresp{
				Status:    0,
				Arid:      arid,
				IsPublish: req.IsPublish,
			}, nil
		}
	}
}

type ESar struct {
	Mogoid     string `json:"mongoid"`
	Fewcontent string `json:"fewcontent"`
	Title      string `json:"title"`
}

func (l *UpdatecontentLogic) ToES(mongoid, title, fewcontent, operation string, isDelete bool) error {
	esclient, e := elastic.NewClient(elastic.SetURL(l.svcCtx.Config.ES.Addr), elastic.SetSniff(false))
	if e != nil {
		return e
	}
	if _, _, er := esclient.Ping(l.svcCtx.Config.ES.Addr).Do(context.Background()); er != nil {
		return er
	}
	tw := ESar{
		Mogoid:     mongoid,
		Fewcontent: fewcontent,
		Title:      title,
	}
	if isDelete {
		_, errrr := esclient.DeleteByQuery().Index("article").Query(elastic.NewTermQuery("mongoid", mongoid)).ProceedOnVersionConflict().Do(context.Background())
		return errrr
	}
	switch operation {
	case "insert":
		if _, err := esclient.Index().Index("article").BodyJson(tw).Do(context.Background()); err != nil {
			return err
		} else {
			return nil
		}
	case "update":
		//!
		script := elastic.NewScriptInline("ctx._source.title=params.title;ctx._source.fewcontent=params.fewcontent").Params(map[string]interface{}{
			"title":      title,
			"fewcontent": fewcontent,
		})
		_, errr := esclient.UpdateByQuery().Index("article").Query(elastic.NewTermQuery("mongoid", mongoid)).Script(script).ProceedOnVersionConflict().Do(context.Background())
		return errr

	default:
		return errors.New("wrong operation")
	}
}

//The first error id for mongoDB the second error is for mysql
func (l *UpdatecontentLogic) ToMongo(ar bson.M, operation string, req *types.Articlereq) (error, error) {
	clientops := options.Client().ApplyURI(l.svcCtx.Config.Mongo.Addr)
	mongoclient, err := mongo.Connect(context.TODO(), clientops)
	if err != nil {
		return err, nil
	}
	collection := mongoclient.Database("DB").Collection("article")
	switch operation {
	case "insert":
		if ar["isDelete"] == true {
			_, errr := collection.DeleteOne(context.Background(), bson.M{"arid": ar["arid"]}, options.Delete().SetCollation(&options.Collation{
				CaseLevel: false,
			}))
			l.ToUserMongo(ar["authorid"].(string), ar["arid"].(string), "delete")
			e := l.svcCtx.ArticleMysqlModel.Delete(l.ctx, ar["arid"].(string))
			return errr, e
		} else {
			_, er := collection.InsertOne(context.Background(), ar)
			l.ToUserMongo(ar["authorid"].(string), ar["arid"].(string), "insert")
			if req.IsPublish  {
				_, e := l.svcCtx.ArticleMysqlModel.Insert(l.ctx, &model.Articles{
					Mongoid:    ar["arid"].(string),
					Title:      ar["title"].(string),
					Fewcontent: ar["fewcontent"].(string),
					Likes:      ar["likes"].(int64),
					Views:      ar["views"].(int64),
					Url:        ar["url"].(string),
					Pubtime:    strconv.Itoa(int(ar["created"].(int64))),
					Coverlinks: ar["coverlink"].(string),
				})
				return er, e
			} else {
				return er, nil
			}
		}
	case "update":
		if ar["isDelete"] == true {
			_, errr := collection.DeleteOne(context.Background(), bson.M{"arid": ar["arid"]}, options.Delete().SetCollation(&options.Collation{
				CaseLevel: false,
			}))
			l.ToUserMongo(ar["authorid"].(string), ar["arid"].(string), "delete")
			e := l.svcCtx.ArticleMysqlModel.Delete(l.ctx, ar["arid"].(string))
			return errr, e
		} else {
			//!
			_, e := collection.UpdateOne(context.Background(), bson.M{"arid": ar["arid"]}, bson.M{"$set": bson.M{"content": ar["content"], "created": ar["created"], "ispublish": ar["ispublish"], "lastrefresh": ar["lastrefresh"], "fewcontent": ar["fewcontent"], "isDelete": ar["isDelete"]}}, options.Update().SetUpsert(true))
			er := l.svcCtx.ArticleMysqlModel.Update(l.ctx, &model.Articles{
				Mongoid:    ar["arid"].(string),
				Title:      ar["title"].(string),
				Fewcontent: ar["fewcontent"].(string),
				Likes:      ar["likes"].(int64),
				Views:      ar["views"].(int64),
				Url:        ar["url"].(string),
				Pubtime:    strconv.Itoa(int(ar["created"].(int64))),
				Coverlinks: ar["coverlink"].(string),
			})
			return e, er
		}
	default:
		return errors.New("wrong operation"), nil
	}
}
func (l *UpdatecontentLogic) ToUserMongo(uid, arid, operation string) error {
	client, e := mongo.Connect(context.TODO(), options.Client().ApplyURI(l.svcCtx.Config.Mongo.Addr))
	if e != nil {
		return e
	}
	collection := client.Database("DB").Collection("userarticle")
	res := collection.FindOne(context.Background(), bson.M{"uid": uid}, options.FindOne())
	var got bson.M
	er := res.Decode(&got)
	if er != nil {
		return er
	}
	arlist := got["articles"].(bson.A)
	switch operation {
	case "insert":
		arlist = append(arlist, arid)
		_, er := collection.UpdateOne(context.Background(), bson.M{"uid": uid}, bson.M{"$set": bson.M{"articles": arlist}})
		if er != nil {
			return er
		}
	case "delete":
		for i := 0; i < len(arlist); i++ {
			if arlist[i] == arid {
				arlist = append(arlist[:i], arlist[i+1:]...)
			}
		}
		_, er := collection.UpdateOne(context.Background(), bson.M{"uid": uid}, bson.M{"$set": bson.M{"articles": arlist}})
		if er != nil {
			return er
		}
	default:
		return errors.New("wrong operaion")
	}
	return nil
}

//for insert bson the arid should be "",
//for update bson the author should be "",
//The returns are bson, arid, title,fewcontent
func (l *UpdatecontentLogic) BsonMFiller(content, cover, author, arid string, time int64, ispublish, isDelete bool) (bson.M, string, string, string) {
	bm := bson.M{}
	if isDelete {
		bm["arid"] = arid
		bm["isDelete"] = true
		return bm, arid, "", ""
	}
	co := strings.Index(content, "\n")
	title := strings.TrimSpace(strings.Trim(content[:co], "#"))
	if arid == "" {
		node, _ := snowflake.NewNode(1)
		arid = node.Generate().String()
	}
	a := strings.Index(content, "![](")
	b := strings.Index(content, "g)")
	if a == b {
		bm["coverlink"] = cover
	} else {
		bm["coverlink"] = content[a+4 : b+1]
	}
	var min = 100
	if len(content) < 100 {
		min = len(content) - strings.Count(content, "#")
	}
	fewcontent := strings.Trim(content, "#")[1:min]

	bm["fewcontent"] = fewcontent
	bm["title"] = title
	bm["content"] = content
	bm["arid"] = arid
	bm["created"] = time
	bm["authorname"] = author
	bm["authorid"] = l.ctx.Value("uid")
	bm["likes"] = int64(0)
	bm["views"] = int64(0)
	bm["readers"] = []string{}
	bm["ispublish"] = ispublish
	bm["lastrefresh"] = time
	bm["isDelete"] = isDelete
	bm["url"] = l.svcCtx.Config.Url.Url + "/api/reading/content/?id" + arid
	bm["daysdata"] = [7]int{}
	return bm, arid, title, fewcontent
}
