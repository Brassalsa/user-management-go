package api

import (
	"net/http"

	"github.com/Brassalsa/user-management-go/internal/api/middlewares"
	"github.com/Brassalsa/user-management-go/internal/db"
	"github.com/Brassalsa/user-management-go/pkg/helpers"
)

func HandleV1Router(r *http.ServeMux, dbfn *db.Database) {
	type Arr []middlewares.HFunc
	dbCtx := middlewares.GroupMiddlewares[*db.Database]

	userV1Route := helpers.RouteStrCl("/api/v1/users")
	// users routes
	r.HandleFunc(userV1Route("POST", "/login"), dbCtx(dbfn, Arr{HandleLoginUser}))
	r.HandleFunc(userV1Route("POST", "/register"), dbCtx(dbfn, Arr{HandleRegisterUser}))

	// secure routes
	ctx := &middlewares.AuthCtx{
		Dbfn: dbfn,
	}
	gpCtx := middlewares.GroupMiddlewares[*middlewares.AuthCtx]

	r.HandleFunc(userV1Route("GET", "/"), gpCtx(ctx, Arr{middlewares.VerifyToken, HandleCheckCurrentUser}))
	r.HandleFunc(userV1Route("POST", "/avatar"), gpCtx(ctx, Arr{middlewares.VerifyToken, HandleUploadAvatar}))
	r.HandleFunc(userV1Route("PUT", "/"), gpCtx(ctx, Arr{middlewares.VerifyToken, HandleUpdateUser}))
	r.HandleFunc(userV1Route("POST", "/password"), gpCtx(ctx, Arr{middlewares.VerifyToken, HandleUpdatePassword}))
	r.HandleFunc(userV1Route("POST", "/delete"), gpCtx(ctx, Arr{middlewares.VerifyToken, HandleDeleteUser}))

	adminV1Route := helpers.RouteStrCl("/api/v1/admin")
	r.HandleFunc(adminV1Route("GET", "/all-users"), gpCtx(ctx, Arr{middlewares.VerifyToken, middlewares.VerifyAdmin, HandleGetAllUsers}))
	r.HandleFunc(adminV1Route("POST", "/user/register"), gpCtx(ctx, Arr{middlewares.VerifyToken, middlewares.VerifyAdmin, HandleRegisterAdmin}))
	r.HandleFunc(adminV1Route("GET", "/user/{id}"), gpCtx(ctx, Arr{middlewares.VerifyToken, middlewares.VerifyAdmin, HandleGetUserById}))
	r.HandleFunc(adminV1Route("DELETE", "/user/{id}"), gpCtx(ctx, Arr{middlewares.VerifyToken, middlewares.VerifyAdmin, HandleDeleteUserByAdmin}))
	r.HandleFunc(adminV1Route("PUT", "/user/{id}"), gpCtx(ctx, Arr{middlewares.VerifyToken, middlewares.VerifyAdmin, HandleModifyUser}))
}
