package settingava

import (
	"net/http"
	"strings"

	"yourbackend/internal/logic/settingava"
	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SettingavaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Settingavareq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := settingava.NewSettingavaLogic(r.Context(), svcCtx)
		//get the img files
		_, uploadfiles, e := r.FormFile("img")
		if e != nil {
			httpx.ErrorCtx(r.Context(), w, e)
		}
		if strings.HasSuffix(uploadfiles.Filename, ".jpg") || strings.HasSuffix(uploadfiles.Filename, ".png") {
			l.Img = uploadfiles
		}
		resp, err := l.Settingava(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}

	}
}
