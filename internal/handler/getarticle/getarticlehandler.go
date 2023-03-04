package getarticle

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yourbackend/internal/logic/getarticle"
	"yourbackend/internal/svc"
	"yourbackend/internal/types"
)

func GetarticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Getarticlereq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := getarticle.NewGetarticleLogic(r.Context(), svcCtx)
		resp, err := l.Getarticle(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
