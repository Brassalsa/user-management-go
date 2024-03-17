package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Brassalsa/user-management-go/internal"
	"github.com/Brassalsa/user-management-go/internal/api/middlewares"
	"github.com/Brassalsa/user-management-go/internal/api/validators"
	"github.com/Brassalsa/user-management-go/internal/db"
	"github.com/Brassalsa/user-management-go/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson"
)

// register user
// admin users can create other admins as well as users
func HandleRegisterAdmin(w http.ResponseWriter, r *http.Request, _ *func()) {
	decoder := json.NewDecoder(r.Body)
	params := db.UserRegister{}
	if err := decoder.Decode(&params); err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprint("Error in parsing json: ", err.Error()))
		return
	}

	if isEmpty := helpers.CheckEmptyStrings([]string{params.Username, params.Password, params.Email, params.Role}); isEmpty {
		helpers.RespondWithError(w, 400, "email, username, password, role are required")
		return
	}

	if err := validators.CheckValidPassword(params.Password); err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}

	if err := validators.CheckValidRole(params.Role); err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}

	if err := internal.HashString(&params.Password); err != nil {
		helpers.RespondWithError(w, 500, "Something went wrong")
		return
	}
	ctx, ok := r.Context().Value(middlewares.CtxKey).(*middlewares.AuthCtx)
	if !ok {
		helpers.RespondWithError(w, 500, "something went wrong")
		return
	}

	// inserting data
	if _, err := ctx.Dbfn.InsertIntoCollection("users", params); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			helpers.RespondWithError(w, 400, "username or email already exists")
			return
		}
		helpers.RespondWithError(w, 500, fmt.Sprint("Registering failed: ", err.Error()))
		return
	}

	helpers.RespondWithJSON(w, 201, "User created successfully!")
}

// delete user by admin
func HandleDeleteUserByAdmin(w http.ResponseWriter, r *http.Request, _ *func()) {
	userID := r.PathValue("id")

	ctx, ok := r.Context().Value(middlewares.CtxKey).(*middlewares.AuthCtx)
	if !ok {
		helpers.RespondWithError(w, 500, "something went wrong")
		return
	}

	objId, err := ctx.Dbfn.ConvertStrToId(userID)
	if err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}

	if ctx.AuthUser.ID == objId {
		helpers.RespondWithError(w, 400, "cann't delete account with this route")
		return
	}

	res, err := ctx.Dbfn.FindOne("users", bson.D{{
		Key:   "_id",
		Value: objId,
	}})

	if err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}

	user := db.User{}
	if err = res.Decode(&user); err != nil {

		if strings.Contains(err.Error(), "no documents in result") {
			helpers.RespondWithError(w, 404, "not found")
			return
		}

		helpers.RespondWithError(w, 400, err.Error())
		return
	}

	if isAdmin := user.Role == "admin"; isAdmin {
		helpers.RespondWithError(w, 400, "cannot delete other admin accounts")
		return
	}

	if err = ctx.Dbfn.DeleteFromCollection("users", bson.D{{
		Key:   "_id",
		Value: user.ID,
	}}); err != nil {
		helpers.RespondWithError(w, 400, err.Error())
	}

	helpers.RespondWithJSON(w, 200, "ok")
}
