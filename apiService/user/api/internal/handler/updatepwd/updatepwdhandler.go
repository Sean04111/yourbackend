package updatepwd

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yourbackend/apiService/user/api/internal/logic/updatepwd"
	"yourbackend/apiService/user/api/internal/svc"
	"yourbackend/apiService/user/api/internal/types"
)

func UpdatepwdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Updatepwdreq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := updatepwd.NewUpdatepwdLogic(r.Context(), svcCtx)
		resp, err := l.Updatepwd(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
