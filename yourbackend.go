package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"

	"yourbackend/internal/config"
	"yourbackend/internal/handler"
	"yourbackend/internal/svc"

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
		c.Static("/", "./internal/static/ava")
		c.Run(":9090")
	}()
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()

}
