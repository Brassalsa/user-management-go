package api

import (
	"net/http"

	"github.com/Brassalsa/user-management-go/internal/controllers"
	"github.com/Brassalsa/user-management-go/pkg/helpers"
)

func HandleV1Router(r *http.ServeMux){
	routeStrCls := helpers.RouteStrClosure("/api/v1")
	// users routes
	r.HandleFunc(routeStrCls("POST", "/users/login"), controllers.HandleLoginUser)
	r.HandleFunc(routeStrCls("POST", "/users/register"), controllers.HandleRegisterUser)
}