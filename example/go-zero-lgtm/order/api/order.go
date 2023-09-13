package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zxmfke/lgtm/example/go-zero-lgtm/order/api/internal/config"
	"github.com/zxmfke/lgtm/example/go-zero-lgtm/order/api/internal/handler"
	"github.com/zxmfke/lgtm/example/go-zero-lgtm/order/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/order.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)

	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	logx.Infof("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
