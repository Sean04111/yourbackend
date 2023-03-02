package searcharticle

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yourbackend/internal/logic/searcharticle"
	"yourbackend/internal/svc"
	"yourbackend/internal/types"
)

func SearchbytitleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Searchreq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := searcharticle.NewSearchbytitleLogic(r.Context(), svcCtx)
		resp, err := l.Searchbytitle(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
