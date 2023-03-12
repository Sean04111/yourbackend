package baseinfo

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yourbackend/apiService/user/api/internal/logic/baseinfo"
	"yourbackend/apiService/user/api/internal/svc"
)

func BaseinfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := baseinfo.NewBaseinfoLogic(r.Context(), svcCtx)
		resp, err := l.Baseinfo()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
