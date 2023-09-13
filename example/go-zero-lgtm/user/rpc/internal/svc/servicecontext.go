package svc

import "github.com/zxmfke/lgtm/example/go-zero-lgtm/user/rpc/internal/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
