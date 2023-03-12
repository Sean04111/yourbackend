// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	baseinfo "yourbackend/apiService/user/api/internal/handler/baseinfo"
	login "yourbackend/apiService/user/api/internal/handler/login"
	pubkey "yourbackend/apiService/user/api/internal/handler/pubkey"
	refreshToken "yourbackend/apiService/user/api/internal/handler/refreshToken"
	register "yourbackend/apiService/user/api/internal/handler/register"
	settingava "yourbackend/apiService/user/api/internal/handler/settingava"
	settingbase "yourbackend/apiService/user/api/internal/handler/settingbase"
	updatepwd "yourbackend/apiService/user/api/internal/handler/updatepwd"
	verification_code "yourbackend/apiService/user/api/internal/handler/verification_code"
	"yourbackend/apiService/user/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
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
				Path:    "/api/login",
				Handler: login.LoginHandler(serverCtx),
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
				Method:  http.MethodPost,
				Path:    "/api/refreshToken",
				Handler: refreshToken.RefreshTokenHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/reguser",
				Handler: register.RegisterHandler(serverCtx),
			},
		},
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
				Path:    "/api/updatepwd",
				Handler: updatepwd.UpdatepwdHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/code",
				Handler: verification_code.VericodeHandler(serverCtx),
			},
		},
	)
}