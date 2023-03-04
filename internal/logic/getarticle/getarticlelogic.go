package getarticle

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"yourbackend/internal/svc"
	"yourbackend/internal/types"
)

type GetarticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetarticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetarticleLogic {
	return &GetarticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetarticleLogic) Getarticle(req *types.Getarticlereq) (*types.Getarticleresp, error) {
	gots, er := l.FromMysql()
	if er != nil {
		return &types.Getarticleresp{
			Status: 1,
		}, nil
	}
	var res []types.GetArticle
	for _, k := range gots {
		res = append(res, types.GetArticle{Id: k.Mongoid, Title: k.Title, Url: k.Url, Description: k.Fewcontent, Likes: k.Likes, Reads: k.Views, Pubtime: k.Pubtime, Imglink: k.Coverlinks})
	}
	return &types.Getarticleresp{
		Status:      0,
		Articlelist: res,
	}, nil
}

type arti struct {
	Mongoid    string
	Title      string
	Fewcontent string
	Likes      int
	Views      int
	Url        string
	Pubtime    string
	Coverlinks string
}

func (l *GetarticleLogic) FromMysql() ([]*arti, error) {
	dbgorm, er := gorm.Open(mysql.Open(l.svcCtx.Config.DB.Mysql), &gorm.Config{})
	if er != nil {
		return nil, er
	}
	db := dbgorm.Table("articles")
	var artis []*arti
	db.Find(&artis)
	return artis, nil
}
