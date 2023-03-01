package updatecontent

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yourbackend/internal/logic/updatecontent"
	"yourbackend/internal/svc"
	"yourbackend/internal/types"
)

func UpdatecontentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Articlereq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := updatecontent.NewUpdatecontentLogic(r.Context(), svcCtx)
		resp, err := l.Updatecontent(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
