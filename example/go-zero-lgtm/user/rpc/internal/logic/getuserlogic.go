package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zxmfke/lgtm/example/go-zero-lgtm/user/rpc/internal/svc"
	"github.com/zxmfke/lgtm/example/go-zero-lgtm/user/rpc/user"
	"github.com/zxmfke/lgtm/example/go-zero-lgtm/user/rpc/userclient"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *userclient.IdRequest) (*userclient.UserResponse, error) {
	return &user.UserResponse{
		Id:   "1",
		Name: "test",
	}, nil
}
