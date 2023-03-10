// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	AI "yourbackend/internal/handler/AI"
	articleread "yourbackend/internal/handler/articleread"
	baseinfo "yourbackend/internal/handler/baseinfo"
	getardata "yourbackend/internal/handler/getardata"
	getarticle "yourbackend/internal/handler/getarticle"
	getbar "yourbackend/internal/handler/getbar"
	getdraft "yourbackend/internal/handler/getdraft"
	getmylikes "yourbackend/internal/handler/getmylikes"
	getsingledata "yourbackend/internal/handler/getsingledata"
	infotable "yourbackend/internal/handler/infotable"
	likearticle "yourbackend/internal/handler/likearticle"
	login "yourbackend/internal/handler/login"
	myarlist "yourbackend/internal/handler/myarlist"
	pubkey "yourbackend/internal/handler/pubkey"
	refreshToken "yourbackend/internal/handler/refreshToken"
	register "yourbackend/internal/handler/register"
	searcharticle "yourbackend/internal/handler/searcharticle"
	settingava "yourbackend/internal/handler/settingava"
	settingbase "yourbackend/internal/handler/settingbase"
	tablename "yourbackend/internal/handler/tablename"
	tools "yourbackend/internal/handler/tools"
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
				Path:    "/api/reguser",
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
				Path:    "/my/updata.content",
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

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/my/allardatas",
				Handler: getardata.GetardataHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/ardata",
				Handler: getsingledata.GetsingledataHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/my/ardatalable",
				Handler: tablename.TablenameHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/my/myarlist",
				Handler: myarlist.MyarlistHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/reading/like",
				Handler: likearticle.LikearticleHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/reading/likeinfo",
				Handler: getbar.GetbarHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/article/getTools",
				Handler: tools.ToolsHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/user/myLikes",
				Handler: getmylikes.GetmylikesHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/user/infoTabs",
				Handler: infotable.InfotableHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/my/draft",
				Handler: getdraft.GetdraftHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/my/chat",
				Handler: AI.AIHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)
}
