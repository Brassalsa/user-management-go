package middlewares

import (
	"net/http"

	"github.com/Brassalsa/user-management-go/pkg/helpers"
)

func VerifyAdmin(w http.ResponseWriter, r *http.Request, mCtx *AuthCtx, next *func()) {
	if isAdmin := mCtx.AuthUser.Role == "admin"; !isAdmin {
		helpers.RespondWithError(w, 401, "unauthorized")
	}

	(*next)()
}
