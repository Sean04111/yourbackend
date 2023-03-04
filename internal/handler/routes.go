// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	articleread "yourbackend/internal/handler/articleread"
	baseinfo "yourbackend/internal/handler/baseinfo"
	getarticle "yourbackend/internal/handler/getarticle"
	login "yourbackend/internal/handler/login"
	pubkey "yourbackend/internal/handler/pubkey"
	refreshToken "yourbackend/internal/handler/refreshToken"
	register "yourbackend/internal/handler/register"
	searcharticle "yourbackend/internal/handler/searcharticle"
	settingava "yourbackend/internal/handler/settingava"
	settingbase "yourbackend/internal/handler/settingbase"
	updatecontent "yourbackend/internal/handler/updatecontent"
	updatepwd "yourbackend/internal/handler/updatepwd"
	verification_code "yourbackend/internal/handler/verification_code"
	"yourbackend/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/code",
				Handler: verification_code.VericodeHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/register",
				Handler: register.RegisterHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/login",
				Handler: login.LoginHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/updatepwd",
				Handler: updatepwd.UpdatepwdHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/refreshToken",
				Handler: refreshToken.RefreshTokenHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/pubkey",
				Handler: pubkey.PubkeyHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/user/baseInfo",
				Handler: baseinfo.BaseinfoHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/myinfo/settingbase",
				Handler: settingbase.SettingbaseHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/myinfo/settingnewimg",
				Handler: settingava.SettingavaHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/updatecontent",
				Handler: updatecontent.UpdatecontentHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/article/searchByTitle",
				Handler: searcharticle.SearchbytitleHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/article/getarticle",
				Handler: getarticle.GetarticleHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/reading/content",
				Handler: articleread.ArticlereadHandler(serverCtx),
			},
		},
	)
}
