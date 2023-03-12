package pubkey

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yourbackend/apiService/user/api/internal/logic/pubkey"
	"yourbackend/apiService/user/api/internal/svc"
)

func PubkeyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := pubkey.NewPubkeyLogic(r.Context(), svcCtx)
		resp, err := l.Pubkey()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
