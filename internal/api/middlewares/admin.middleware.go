package middlewares

import (
	"net/http"

	"github.com/Brassalsa/user-management-go/pkg/helpers"
)

func VerifyAdmin(w http.ResponseWriter, r *http.Request, next *func()) {

	ctx, ok := r.Context().Value(CtxKey).(*AuthCtx)
	if !ok {
		helpers.RespondWithError(w, 500, "something went wrong")
		return
	}

	// check if admin
	if isAdmin := ctx.AuthUser.Role == "admin"; !isAdmin {
		helpers.RespondWithError(w, 401, "unauthorized")
		return
	}

	(*next)()
}
