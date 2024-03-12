package controllers

import (
	"net/http"

	"github.com/Brassalsa/user-management-go/pkg/helpers"
)

func LoginUser(w http.ResponseWriter, _ *http.Request) {

	helpers.RespondWithJSON(w, 200, "User created")
}

func RegisterUser(w http.ResponseWriter, _ *http.Request) {
	helpers.RespondWithJSON(w, 200, "User created")
}