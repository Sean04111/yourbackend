package register

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"time"

	"yourbackend/internal/model"
	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/bwmarrin/snowflake"
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
	//get the code for user
	usercode, e := l.FromRedis(req.Email)
	if e != nil {
		return &types.Registerresp{
			Status: 1,
		}, nil
	}
	//user code not right
	if usercode != req.Code {
		return &types.Registerresp{
			Status: 2,
		}, nil
	} else { //user code right
		checkcode, b := l.FromRedis(usercode)
		if b != nil {
			return &types.Registerresp{
				Status: 1,
			}, nil
		}
		if checkcode == req.Check {
			node, er := snowflake.NewNode(1)
			if er != nil {
				return &types.Registerresp{
					Status: 1,
				}, nil
			}
			uid := node.Generate()
			//decode the password
			stringpwd, e := l.svcCtx.RsaOps.Decode(req.Pass)
			if e != nil {
				return &types.Registerresp{
					Status: 1,
				}, nil
			}
			storepwd, erro := bcrypt.GenerateFromPassword(stringpwd, 10)
			if erro != nil {
				return &types.Registerresp{
					Status: 1,
				}, erro
			}
			newuser := model.User{
				Uid:      uid.String(),
				Email:    req.Email,
				Password: string(storepwd),
				Name:     req.Name,
			}
			_, errorr := l.svcCtx.MysqlModel.Insert(l.ctx, &newuser)
			if errorr != nil {
				log.Fatalln("[Model] Failed to insert data", e)
				return &types.Registerresp{
					Status: 1,
				}, e
			}
			now := time.Now().Unix()
			Jwttoken, err := l.GetJWT(l.svcCtx.Config.Auth.AccessSecret, newuser.Email, newuser.Uid, l.svcCtx.Config.Auth.AccessExpire, now)
			if err != nil {
				log.Fatalln("[JWT] Failed to generate json web token : ", err)
				return &types.Registerresp{
					Status: 1,
				}, err
			}
			errr := l.ToMongo(uid.String())
			if errr != nil {
				return &types.Registerresp{
					Status: 1,
				}, nil
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
		Addr: l.svcCtx.Config.Redis.Addr,
		DB:   l.svcCtx.Config.Redis.DB,
	})
	code, e := redisclient.Get(keyword).Result()
	if e != nil {
		return "", e
	}
	return code, nil
}
func (l *RegisterLogic) GetJWT(secretkey, email, uid string, lasttime, starttime int64) (string, error) {
	claim := make(jwt.MapClaims)
	claim["starttime"] = starttime
	claim["uid"] = uid
	claim["expiretime"] = starttime + lasttime
	claim["email"] = email
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claim
	return token.SignedString([]byte(secretkey))
}
func (l *RegisterLogic) ToMongo(uid string) error {
	mongoclient, e := mongo.Connect(context.Background(), options.Client().ApplyURI(l.svcCtx.Config.Mongo.Addr))
	if e != nil {
		return e
	}
	collection := mongoclient.Database("DB").Collection("userarticle")
	toinsert:=bson.M{}
	toinsert["uid"] = uid
	toinsert["articles"] = []string{}
	toinsert["alldata"] = [7]int{}
	toinsert["lastrefresh"] = time.Now().Unix()
	toinsert["likes"] = []string{}
	_, err := collection.InsertOne(context.Background(), toinsert)
	if err != nil {
		return err
	}
	return nil
}
