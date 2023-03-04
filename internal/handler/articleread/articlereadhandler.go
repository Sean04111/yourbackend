package articleread

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yourbackend/internal/logic/articleread"
	"yourbackend/internal/svc"
	"yourbackend/internal/types"
)

func ArticlereadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Readarticlereq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := articleread.NewArticlereadLogic(r.Context(), svcCtx)
		resp, err := l.Articleread(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
