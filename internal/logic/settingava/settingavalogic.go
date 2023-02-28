package settingava

import (
	"context"
	"mime/multipart"
	"os"

	"yourbackend/internal/model"
	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SettingavaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	Img *multipart.FileHeader
}

func NewSettingavaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SettingavaLogic {
	return &SettingavaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SettingavaLogic) Settingava(req *types.Settingavareq) ( *types.Settingavaresp,  error) {
	newimf,er:=os.Create("internal/static/ava/"+l.ctx.Value("uid").(string)+".jpg")
	defer func ()  {
		if e1:=newimf.Close();e1!=nil{
			panic(e1)
		}	
	}()
	if er!=nil{
		return &types.Settingavaresp{
			Status: 1,
		},nil
	}
	file,e:=l.Img.Open()
	defer func() {
		if e2:=file.Close();e2!=nil{
			panic(e2)
		}
	}()
	if e!=nil{
		return &types.Settingavaresp{
			Status: 1,
		},nil
	}
	date:=make([]byte,1000000)
	_,a:=file.Read(date)
	if a!=nil{
		return &types.Settingavaresp{
			Status: 1,
		},nil
	}
	_,b:=newimf.Write(date)
	if b!=nil{
		return &types.Settingavaresp{
			Status: 1,
		},nil
	}
	link:="127.0.0.1:8888/internal/ava/"+l.ctx.Value("uid").(string)+".jpg"//The static router needed!
	gotuser,err:=l.svcCtx.MysqlModel.FindOne(l.ctx,l.ctx.Value("email").(string))
	if err!=nil{
		return &types.Settingavaresp{
			Status: 1,
		},nil
	}
	if l.svcCtx.MysqlModel.Update(l.ctx,&model.User{
		Uid: gotuser.Uid,
		Email: gotuser.Email,
		Password: gotuser.Password,
		Name:gotuser.Name,
		AvatarLink: link,
		Profession: gotuser.Profession,
		Type:gotuser.Type,
	})!=nil{
		return &types.Settingavaresp{
			Status :1,
		},nil
	}
	return &types.Settingavaresp{
		Status: 0,
	},nil
}


