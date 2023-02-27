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
	"golang.org/x/crypto/bcrypt"
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
	usercode, e := l.FromRedis(req.Email)
	if e != nil {
		return &types.Registerresp{
			Status: 1,
		}, nil
	}
	if usercode != req.Code {
		return &types.Registerresp{
			Status: 2,
		}, nil
	} else {
		checkcode, b := l.FromRedis(usercode)
		if b != nil {
			return &types.Registerresp{
				Status: 1,
			}, nil
		}
		if checkcode == req.Check {
			stringuid, er := l.FromRedis("usernum")
			if er != nil {
				log.Fatalln("[Redis] Failed to get the data:")
				return &types.Registerresp{
					Status: 1,
				}, nil
			}
			uid, _ := strconv.Atoi(stringuid)
			storepwd, erro := bcrypt.GenerateFromPassword(l.svcCtx.RsaOps.Decode([]byte(req.Pass)), 10)
			if erro != nil {
				return &types.Registerresp{
					Status: 1,
				}, erro
			}
			newuser := model.User{
				Uid:      int64(uid) + 1,
				Email:    req.Email,
				Password: string(storepwd),
				Name:     req.Name,
			}
			_, e := l.svcCtx.MysqlModel.Insert(l.ctx, &newuser)
			if e != nil {
				log.Fatalln("[Model] Failed to insert data", e)
				return &types.Registerresp{
					Status: 1,
				}, e
			}
			now := time.Now().Unix()
			Jwttoken, err := l.GetJWT(l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire, newuser.Uid, now)
			if err != nil {
				log.Fatalln("[JWT] Failed to generate json web token : ", err)
				return &types.Registerresp{
					Status: 1,
				}, err
			}
			return &types.Registerresp{
				Status:      0,
				Accesstoken: Jwttoken,
				Expires:     strconv.Itoa(int(l.svcCtx.Config.Auth.AccessExpire + now)), //too long !
				Name:        req.Name,
			}, nil
		} else {
			return &types.Registerresp{
				Status: 3,
			}, nil

		}
	}
}
func (l *RegisterLogic) FromRedis(keyword string) (string, error) {
	redisclient := redis.NewClient(&redis.Options{
		Addr: "121.36.131.50:6379",
		DB:   0,
	})
	if code, e := redisclient.Get(keyword).Result(); e != nil {
		log.Fatalln("[Redis] Failed to get data from redis")
		return "", e
	} else {
		return code, nil
	}
}
func (l *RegisterLogic) GetJWT(secretkey string, lasttime, uid, starttime int64) (string, error) {
	claim := make(jwt.MapClaims)
	claim["starttime"] = starttime
	claim["uid"] = uid
	claim["expiretime"] = starttime + lasttime
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claim
	return token.SignedString([]byte(secretkey))
}