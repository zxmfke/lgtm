package handler

import (
	"github.com/zxmfke/lgtm/example/go-zero-lgtm/order/api/internal/logic"
	"github.com/zxmfke/lgtm/example/go-zero-lgtm/order/api/internal/svc"
	"github.com/zxmfke/lgtm/example/go-zero-lgtm/order/api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getOrderHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetOrderLogic(r.Context(), ctx)
		resp, err := l.GetOrder(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
