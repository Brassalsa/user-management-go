package api

import (
	"net/http"

	"github.com/Brassalsa/user-management-go/pkg/helpers"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	helpers.RespondWithJSON(w, 200, "Its live")
}

func HandlerError(w http.ResponseWriter, r *http.Request) {
	helpers.RespondWithError(w, 400, "Something went wrong")
}
