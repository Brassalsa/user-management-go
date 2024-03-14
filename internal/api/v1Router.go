package api

import (
	"net/http"

	"github.com/Brassalsa/user-management-go/internal/api/controllers"
	"github.com/Brassalsa/user-management-go/internal/api/middlewares"
	"github.com/Brassalsa/user-management-go/internal/db"

	"github.com/Brassalsa/user-management-go/pkg/helpers"
)

func HandleV1Router(r *http.ServeMux, dbfn *db.Database) {
	dbCtx := ApiContext[*db.Database]{
		ctx: dbfn,
	}
	routeStrCls := helpers.RouteStrClosure("/api/v1")
	// users routes
	r.HandleFunc(routeStrCls("POST", "/users/login"), dbCtx.Provider(controllers.HandleLoginUser))
	r.HandleFunc(routeStrCls("POST", "/users/register"), dbCtx.Provider(controllers.HandleRegisterUser))

	// secure routes
	mdCtx := ApiContext[*middlewares.AuthCtx]{
		ctx: &middlewares.AuthCtx{
			Dbfn: dbfn,
		},
	}
	gpCtx := middlewares.GroupMiddlewares[*middlewares.AuthCtx]
	type Arr []middlewares.HFunc[*middlewares.AuthCtx]
	r.HandleFunc(routeStrCls("POST", "/users/current"), gpCtx(mdCtx.ctx, Arr{middlewares.VerifyTokenfunc, controllers.HandleCheckCurrentUser}))
}
