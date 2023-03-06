package likearticle

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"yourbackend/internal/model"

	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikearticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikearticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikearticleLogic {
	return &LikearticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikearticleLogic) Likearticle(req *types.Likearticlereq) (resp *types.Likearticleresp, err error) {
	mongoclient, e := mongo.Connect(context.Background(), options.Client().ApplyURI(l.svcCtx.Config.Mongo.Addr))
	if e != nil {
		return &types.Likearticleresp{
			Status: 1,
		}, nil
	}
	r := mongoclient.Database("DB").Collection("user-article")
	res := r.FindOne(context.Background(), bson.M{"uid": l.ctx.Value("uid").(string)})
	var got bson.M
	errr := res.Decode(&got)
	if errr != nil {
		return &types.Likearticleresp{
			Status: 1,
		}, nil
	}
	gotlikes := got["likes"].(bson.A)
	if req.Active {
		gotlikes = append(gotlikes, req.Id)
	} else {
		for i, k := range gotlikes {
			if k.(string) == req.Id {
				gotlikes = append(gotlikes[:i], gotlikes[i+1:])
			}
		}
	}
	_, er := r.UpdateOne(context.Background(), bson.M{"uid": l.ctx.Value("uid").(string)}, bson.M{"$set": bson.M{"likes": gotlikes}})
	if er != nil {
		return &types.Likearticleresp{
			Status: 1,
		}, nil
	}
	errrr := l.Datahandler(req)
	if errrr != nil {
		return &types.Likearticleresp{
			Status: 1,
		}, nil
	}
	return &types.Likearticleresp{
		Status: 0,
	}, nil
}
func (l *LikearticleLogic) Datahandler(req *types.Likearticlereq) error {
	mongoclient, e := mongo.Connect(context.Background(), options.Client().ApplyURI(l.svcCtx.Config.Mongo.Addr))
	if e != nil {
		return e
	}
	r := mongoclient.Database("DB").Collection("article")
	res := r.FindOne(context.Background(), bson.M{"arid": req.Id})
	var got bson.M
	res.Decode(&got)
	if req.Active {
		_, err := r.UpdateOne(context.Background(), bson.M{"arid": req.Id}, bson.M{"$set": bson.M{"likes": got["likes"].(int64) + 1}})
		if err != nil {
			return err
		}
		errr := l.svcCtx.ArticleMysqlModel.Update(l.ctx, &model.Articles{
			Mongoid:    req.Id,
			Title:      got["title"].(string),
			Fewcontent: got["fewcontent"].(string),
			Likes:      got["likes"].(int64) + 1,
			Views:      got["views"].(int64),
			Url:        got["url"].(string),
			Pubtime:    got["created"].(string),
			Coverlinks: got["coverlink"].(string),
		})
		if errr != nil {
			return errr
		}
		return nil
	} else {
		_, err := r.UpdateOne(context.Background(), bson.M{"arid": req.Id}, bson.M{"$set": bson.M{"likes": got["likes"].(int64) - 1}})
		if err != nil {
			return err
		}
		errr := l.svcCtx.ArticleMysqlModel.Update(l.ctx, &model.Articles{
			Mongoid:    req.Id,
			Title:      got["title"].(string),
			Fewcontent: got["fewcontent"].(string),
			Likes:      got["likes"].(int64) - 1,
			Views:      got["views"].(int64),
			Url:        got["url"].(string),
			Pubtime:    got["created"].(string),
			Coverlinks: got["coverlink"].(string),
		})
		if errr != nil {
			return errr
		}
		return nil
	}
}
