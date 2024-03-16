package api

import (
	"net/http"

	"github.com/Brassalsa/user-management-go/internal/api/middlewares"
	"github.com/Brassalsa/user-management-go/internal/db"
	"github.com/Brassalsa/user-management-go/pkg/helpers"
)

func HandleV1Router(r *http.ServeMux, dbfn *db.Database) {
	dbCtx := ApiContext[*db.Database]{
		ctx: dbfn,
	}
	v1Route := helpers.RouteStrClosure("/api/v1")
	// users routes
	r.HandleFunc(v1Route("POST", "/users/login"), dbCtx.Provider(HandleLoginUser))
	r.HandleFunc(v1Route("POST", "/users/register"), dbCtx.Provider(HandleRegisterUser))

	// secure routes
	mdCtx := ApiContext[*middlewares.AuthCtx]{
		ctx: &middlewares.AuthCtx{
			Dbfn: dbfn,
		},
	}
	gpCtx := middlewares.GroupMiddlewares[*middlewares.AuthCtx]
	type Arr []middlewares.HFunc[*middlewares.AuthCtx]
	r.HandleFunc(v1Route("GET", "/users/me"), gpCtx(mdCtx.ctx, Arr{middlewares.VerifyTokenfunc, HandleCheckCurrentUser}))
	r.Handle(v1Route("POST", "/users/avatar"), gpCtx(mdCtx.ctx, Arr{middlewares.VerifyTokenfunc, HandleUploadAvatar}))

}
