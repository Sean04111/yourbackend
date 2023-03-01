package updatecontent

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/google/uuid"
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
//redis needed as a cache
func (l *UpdatecontentLogic) Updatecontent(req *types.Articlereq) (*types.Articleresp, error) {
	if req.Arid == "" {
		gotuser, e := l.svcCtx.MysqlModel.FindOne(l.ctx, l.ctx.Value("email").(string))
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
			err := l.ToES(arid, title, fewcontent, "insert")
			if err != nil {
				errchan <- err
			} else {
				return
			}
		}()
		go func() {
			defer wg.Done()
			er := l.ToMongo(bson, "insert")
			if er != nil {
				errchan <- er
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
				Status:    1,
				Arid:      arid,
				IsPublish: req.IsPublish,
			}, nil
		}
	} else {
		bson, arid, title, fewcontent := l.BsonMFiller(req.Content, "", "", req.Arid, time.Now().Unix(), req.IsPublish, req.Content == "")
		errchan := make(chan error, 3)
		wg := new(sync.WaitGroup)
		wg.Add(2)
		go func() {
			defer wg.Done()
			err := l.ToMongo(bson, "update")
			if err != nil {
				errchan<-err
			}else{
				return
			}
		}()
		go func(){
			defer wg.Done()
			err:=l.ToES(arid,title,fewcontent,"update")
			if err!=nil{
				errchan<-err
			}else{
				return
			}
		}()
		wg.Wait()
		select{
		case <-errchan:
			return &types.Articleresp{
				Status: 1,
			},nil
			default:
				return &types.Articleresp{
					Status: 0,
					Arid: arid,
					IsPublish: req.IsPublish,
				},nil
		}
	}
}

type ESar struct {
	Mogoid     string `json:"mongoid"`
	Fewcontent string `json:"fewcontent"`
	Title      string `json:"title"`
}

func (l *UpdatecontentLogic) ToES(mongoid, title, fewcontent, operation string) error {
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
	switch operation {
	case "insert":
		if _, err := esclient.Index().Index("article").BodyJson(tw).Do(context.Background()); err != nil {
			return err
		} else {
			return nil
		}
	case "update":
		//!
		script := "ctx._source[`title`]=" + title + ";ctx._source[`fewcontent`]=" + fewcontent 
		_, errr := esclient.UpdateByQuery().Index("article").Query(elastic.NewTermQuery("mongoid", mongoid)).Script(elastic.NewScript(script)).ProceedOnVersionConflict().Do(context.Background())
		return errr
	case "delete":
		_, errrr := esclient.DeleteByQuery().Index("article").Query(elastic.NewTermQuery("mongoid", mongoid)).ProceedOnVersionConflict().Do(context.Background())
		return errrr
	default:
		return errors.New("wrong operation")
	}
}
func (l *UpdatecontentLogic) ToMongo(ar bson.M, operation string) error {
	clientops := options.Client().ApplyURI(l.svcCtx.Config.Mongo.Addr)
	mongoclient, err := mongo.Connect(context.TODO(), clientops)
	if err != nil {
		return err
	}
	collection := mongoclient.Database("DB").Collection("article")
	switch operation {
	case "insert":
		if ar["isDelete"] == true {
			_, errr := collection.DeleteOne(context.Background(), bson.M{"arid": ar["arid"]}, options.Delete().SetCollation(&options.Collation{
				CaseLevel: false,
			}))
			return errr
		} else {
			_, er := collection.InsertOne(context.Background(), ar)
			return er
		}
	case "update":
		if ar["isDelete"] == true {
			_, errr := collection.DeleteOne(context.Background(), bson.M{"arid": ar["arid"]}, options.Delete().SetCollation(&options.Collation{
				CaseLevel: false,
			}))
			return errr
		} else {
			//!
			_, e := collection.UpdateOne(context.Background(), bson.M{"arid": ar["arid"]}, bson.M{"$set": bson.M{"content": ar["content"], "created": ar["created"], "ispublish": ar["ispublish"], "lastrefresh": ar["lastrefresh"], "fewcontent": ar["fewcontent"], "isDelete": ar["isDelete"]}}, options.Update().SetUpsert(true))
			return e
		}
	default:
		return errors.New("wrong operation")
	}
}

//for insert bson the arid should be "",
//for update bson the author should be "",
//The returns are bson, arid, title,fewcontent
func (l *UpdatecontentLogic) BsonMFiller(content, cover, author, arid string, time int64, ispublish, isdelete bool) (bson.M, string, string, string) {
	bm := bson.M{}
	co := strings.Index(content, "\n")
	title := strings.TrimSpace(strings.Trim(content[:co], "#"))
	if arid == "" {
		arid = uuid.New().String()
	}
	a := strings.Index(content, "![](")
	b := strings.Index(content, "g)")
	coverlink := content[a+4 : b+1]
	fewcontent := strings.Trim(strings.Trim(content, title), "#")[1:100]
	if coverlink == "" {
		bm["coverlink"] = coverlink
	} else {
		bm["coverlink"] = cover
	}
	bm["arid"] = arid
	bm["title"] = title
	bm["content"] = content
	bm["created"] = time
	bm["authorname"] = author
	bm["authorid"] = l.ctx.Value("uid")
	bm["likes"] = 0
	bm["views"] = 0
	bm["readers"] = []string{}
	bm["ispublish"] = ispublish
	bm["lastrefresh"] = int64(0)
	bm["fewcontent"] = fewcontent
	bm["isDelete"] = isdelete
	bm["url"] = l.svcCtx.Config.Url.Url + "reading/?ar_id" + arid
	return bm, arid, title, fewcontent
}