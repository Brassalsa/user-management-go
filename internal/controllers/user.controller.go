package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Brassalsa/user-management-go/pkg/helpers"
)

type LoginUsers struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
func HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := LoginUsers{}
	err := decoder.Decode(&params)
	if err !=nil {
		helpers.RespondWithError(w, 400, fmt.Sprint("Error parsing json: ", err))
		return
	}
	fmt.Print(params)

	helpers.RespondWithJSON(w, 200, "Logged in")
}

func HandleRegisterUser(w http.ResponseWriter, _ *http.Request) {
	helpers.RespondWithJSON(w, 201, "User created successfully!")
}