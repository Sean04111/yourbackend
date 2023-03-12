package verification_code

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yourbackend/apiService/user/api/internal/logic/verification_code"
	"yourbackend/apiService/user/api/internal/svc"
	"yourbackend/apiService/user/api/internal/types"
)

func VericodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Codereq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := verification_code.NewVericodeLogic(r.Context(), svcCtx)
		resp, err := l.Vericode(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
