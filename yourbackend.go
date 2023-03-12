package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"yourbackend/internal/config"
	"yourbackend/internal/handler"
	"yourbackend/internal/svc"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/yourbackend-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	go func() {
		c := gin.Default()
		c.Use(Cors())
		c.Static("/", "D:/GOPATH__MY/src/yourbackend/internal/static/")
		c.Use(static.Serve("/", static.LocalFile("internal/static", true)))
		c.NoRoute(func(c *gin.Context) {
			accept := c.Request.Header.Get("Accept")
			flag := strings.Contains(accept, "application/json")
			if !flag {
				content, err := ioutil.ReadFile("D:/GOPATH__MY/src/yourbackend/internal/static/index.html")
				if err != nil {
					c.Writer.WriteHeader(404)
					return
				}
				c.Writer.WriteHeader(200)
				c.Writer.Header().Add("Accept", "text/html")
				c.Writer.Write(content)
				c.Writer.Flush()
			}
		})
		c.SetTrustedProxies([]string{"127.0.0.1:8888"})
		c.Run(":9090")
	}()
	//staticFileHandler(server)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()

}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("11111111111111")
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

//定义函数
func staticFileHandler(engine *rest.Server) {
	//这里注册
	patern := "/static/"
	dirpath := "./internal/static/"

	rd, _ := ioutil.ReadDir(dirpath)

	//添加进路由最后生成 /asset

	for _, f := range rd {
		filename := f.Name()
		path := "/static/" + filename
		//最后生成 /asset
		engine.AddRoute(
			rest.Route{
				Method:  http.MethodGet,
				Path:    path,
				Handler: dirhandler(patern, dirpath),
			})
	}

}

func dirhandler(patern, filedir string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		handler := http.StripPrefix(patern, http.FileServer(http.Dir(filedir)))
		handler.ServeHTTP(w, req)
	}
}
