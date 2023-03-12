package settingava

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yourbackend/apiService/user/api/internal/logic/settingava"
	"yourbackend/apiService/user/api/internal/svc"
	"yourbackend/apiService/user/api/internal/types"
)

func SettingavaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Settingavareq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := settingava.NewSettingavaLogic(r.Context(), svcCtx)
		resp, err := l.Settingava(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
