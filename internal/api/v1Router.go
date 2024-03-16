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
	userV1Route := helpers.RouteStrClosure("/api/v1/users")
	// users routes
	r.HandleFunc(userV1Route("POST", "/login"), dbCtx.Provider(HandleLoginUser))
	r.HandleFunc(userV1Route("POST", "/register"), dbCtx.Provider(HandleRegisterUser))

	// secure routes
	mdCtx := ApiContext[*middlewares.AuthCtx]{
		ctx: &middlewares.AuthCtx{
			Dbfn: dbfn,
		},
	}
	gpCtx := middlewares.GroupMiddlewares[*middlewares.AuthCtx]
	type Arr []middlewares.HFunc[*middlewares.AuthCtx]
	r.HandleFunc(userV1Route("GET", "/"), gpCtx(mdCtx.ctx, Arr{middlewares.VerifyToken, HandleCheckCurrentUser}))
	r.HandleFunc(userV1Route("POST", "/avatar"), gpCtx(mdCtx.ctx, Arr{middlewares.FilesMiddleware, middlewares.VerifyToken, HandleUploadAvatar}))
	r.HandleFunc(userV1Route("PUT", "/"), gpCtx(mdCtx.ctx, Arr{middlewares.VerifyToken, HandleUpdateUser}))
	r.HandleFunc(userV1Route("POST", "/password"), gpCtx(mdCtx.ctx, Arr{middlewares.VerifyToken, HandleUpdatePassword}))
	r.HandleFunc(userV1Route("POST", "/delete"), gpCtx(mdCtx.ctx, Arr{middlewares.VerifyToken, HandleDeleteUser}))
}
