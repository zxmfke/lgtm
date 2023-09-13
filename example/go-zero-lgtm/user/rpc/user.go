package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zxmfke/lgtm/example/go-zero-lgtm/user/rpc/internal/config"
	"github.com/zxmfke/lgtm/example/go-zero-lgtm/user/rpc/internal/server"
	"github.com/zxmfke/lgtm/example/go-zero-lgtm/user/rpc/internal/svc"

	"github.com/zxmfke/lgtm/example/go-zero-lgtm/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	srv := server.NewUserServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		userclient.RegisterUserServer(grpcServer, srv)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	logx.Infof("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
