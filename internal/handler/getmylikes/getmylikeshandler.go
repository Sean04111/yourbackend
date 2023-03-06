package getmylikes

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yourbackend/internal/logic/getmylikes"
	"yourbackend/internal/svc"
)

func GetmylikesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := getmylikes.NewGetmylikesLogic(r.Context(), svcCtx)
		resp, err := l.Getmylikes()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
