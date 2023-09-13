package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zxmfke/lgtm/example/go-zero-lgtm/order/api/internal/config"
	"github.com/zxmfke/lgtm/example/go-zero-lgtm/user/rpc/user"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc user.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: user.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
