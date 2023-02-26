package register

import (
	"context"
	"log"
	"strconv"
	"time"

	"yourbackend/internal/model"
	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.Registerreq) (resp *types.Registerresp, err error) {
	realcode := l.FromRedis(req.Email)
	if realcode != req.Code {
		return &types.Registerresp{
			Status: 2,
		}, nil
	}else{
		uid,er:=l.FromuserRedis()
		if er!=nil{
			log.Fatalln("[Redis] Failed to get the data:")
			return &types.Registerresp{
				Status: 1,
			},er
		}
		newuser:=model.User{
			Uid: int64(uid)+1,
			Email:req.Email,
			Password: req.Pass,
			Name: req.Name,
		}
		_,e:=l.svcCtx.MysqlModel.Insert(l.ctx,&newuser)
		if e!=nil{
			log.Fatalln("[Model] Failed to insert data",e)
			return &types.Registerresp{
				Status: 1,
			},e
		}
		now:=time.Now().Unix()
		Jwttoken,err:=l.GetJWT(l.svcCtx.Config.Auth.Secretkey,l.svcCtx.Config.Auth.Expiretime,newuser.Uid,now)
		if err!=nil{
			log.Fatalln("[JWT] Failed to generate json web token : ",err)
			return &types.Registerresp{
				Status: 1,
			},err
		}
		return &types.Registerresp{
			Status: 0,
			Accesstoken: Jwttoken,
			Expires: strconv.Itoa(int(l.svcCtx.Config.Auth.Expiretime+now)),//too long !
			Name:req.Name,
		},nil
	}
}
func (l *RegisterLogic) FromRedis(email string) string {
	redisclient := redis.NewClient(&redis.Options{
		Addr: "121.36.131.50:6379",
		DB:   0,
	})
	if pong, err := redisclient.Ping().Result(); err != nil {
		log.Fatalln("[Redis] Failed to connnect to  redis", err)
		return "" //ugly
	} else if pong != "PONG" {
		log.Fatalln("[Redis] Failed to connect to the redis client")
		return ""
	} else {
		if code, e := redisclient.Get(email).Result(); e != nil {
			log.Fatalln("[Redis] Failed to get data from redis")
			return ""
		} else {
			return code
		}
	}
}
func(l *RegisterLogic)FromuserRedis()(int,error){
	redisclient:=redis.NewClient(&redis.Options{
		Addr: "121.36.131.50:6379",
		DB:   0,
	})
	if n,err:=redisclient.Get("usernum").Result();err!=nil{
		log.Fatalln("[Redis] Failed to get the data: ",err)
		return 0,err
	}else{
		num,_:=strconv.Atoi(n)
		newnum:=num+1
		newstringnum:=strconv.Itoa(newnum)
		if err:=redisclient.Set("usernum",newstringnum,0).Err();err!=nil{
			log.Fatalln("[Redis] Failed to insert data into redis")
		}
		return num,nil

	}
}
func (l *RegisterLogic) GetJWT(secretkey string, lasttime, uid ,starttime int64)(string,error){
	claim:=make(jwt.MapClaims)
	claim["starttime"]=starttime
	claim["uid"]=uid
	claim["expiretime"]=starttime+lasttime
	token:=jwt.New(jwt.SigningMethodHS256)
	token.Claims=claim
	return token.SignedString([]byte(secretkey))
}
