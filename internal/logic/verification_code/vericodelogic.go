package verification_code

import (
	"context"
	"log"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"

	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/go-redis/redis"
	"github.com/jordan-wright/email"
	"github.com/zeromicro/go-zero/core/logx"
)

type VericodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVericodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VericodeLogic {
	return &VericodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VericodeLogic) Vericode(req *types.Codereq) (resp *types.Coderesp, err error) {

	code, err := l.SendCode(req.Email)
	rand.Seed(time.Now().UnixNano())
	frontendcode := rand.Intn(10000)
	if err != nil {
		log.Fatalln("[Email] Send error : ", err)
		return &types.Coderesp{
			Status: 1,
		}, err
	}
	if e := l.ToRedis(req.Email, code); e != nil {
		log.Fatalln("[Redis] insert error : ", e)
		return &types.Coderesp{
			Status: 1,
		}, e
	} else {

		if l.ToRedis(code, strconv.Itoa(frontendcode)) != nil {
			return &types.Coderesp{
				Status: 1,
			}, nil
		}
	}
	return &types.Coderesp{
		Status: 0,
		Code:   strconv.Itoa(frontendcode),
	}, nil
}
func (l *VericodeLogic) SendCode(receiver string) (string, error) {
	//To form a veri code
	rand.Seed(time.Now().UnixNano())
	usercode := rand.Intn(10000)
	//Input the code
	text := "Welcome ! your verrification code is(do NOT tell someone else) :" + strconv.Itoa(usercode)
	em := email.NewEmail()
	em.From = "seantown1998@163.com"
	em.To = []string{receiver}
	em.Subject = "Welcome"
	em.Text = []byte(text)
	err := em.Send("smtp.163.com:25", smtp.PlainAuth("", "seantown1998@163.com", "UROSELZDZDTPFSDV", "smtp.163.com"))
	if err != nil {
		log.Fatalln("[Email] send error :", err)
	}
	//input code into the cache
	return strconv.Itoa(usercode), err
}
func (l *VericodeLogic) ToRedis(keyword, code string) error {
	client := redis.NewClient(&redis.Options{
		Addr: "121.36.131.50:6379", //find a avaliable redis!
		DB:   0,
	})
	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatalln(err)
		return err
	}
	if pong != "PONG" {
		log.Fatalln("[Redis] Failed to connect to the redis client")
	} else {
		err := client.Set(keyword, code, time.Second*180).Err()
		if err != nil {
			log.Fatalln("[Redis] Failed to insert the KV")
			return err
		}
	}
	return nil
}
