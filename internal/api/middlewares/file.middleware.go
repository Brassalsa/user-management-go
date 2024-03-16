package middlewares

import (
	"net/http"

	"github.com/Brassalsa/user-management-go/internal/api/constants"
	"github.com/Brassalsa/user-management-go/pkg/helpers"
)

func FilesMiddleware(w http.ResponseWriter, r *http.Request, mCtx *AuthCtx, next *func()) {
	r.Body = http.MaxBytesReader(w, r.Body, constants.MAX_UPLOAD_SIZE)

	if err := r.ParseMultipartForm(constants.MAX_UPLOAD_SIZE); err != nil {
		helpers.RespondWithError(w, 400, "The uploaded file is too big. Please choose an file that's less than 20MB in size")
		return
	}

	(*next)()
}
