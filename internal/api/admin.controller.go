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
func HandleRegisterAdmin(w http.ResponseWriter, r *http.Request, _ middlewares.Next) {
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

// get user details by id
func HandleGetUserById(w http.ResponseWriter, r *http.Request, _ middlewares.Next) {
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

	user := db.UserWithoutPassword{}
	res, err := ctx.Dbfn.FindOne("users", bson.D{
		{Key: "_id", Value: objId},
	})

	if err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}

	if err = res.Decode(&user); err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			helpers.RespondWithError(w, 404, "not found")
			return
		}
		helpers.RespondWithError(w, 400, err.Error())
		return
	}

	helpers.RespondWithJSON(w, 200, user)
}

// delete user by admin
func HandleDeleteUserByAdmin(w http.ResponseWriter, r *http.Request, _ middlewares.Next) {
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

// get users list
func HandleGetAllUsers(w http.ResponseWriter, r *http.Request, _ middlewares.Next) {
	allUsers := []db.UserWithoutPassword{}
	ctx, ok := r.Context().Value(middlewares.CtxKey).(*middlewares.AuthCtx)
	if !ok {
		helpers.RespondWithError(w, 500, "something went wrong")
		return
	}
	if err := ctx.Dbfn.GetAllDocuments("users", &allUsers); err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}

	helpers.RespondWithJSON(w, 200, allUsers)
}

// modify user details
type modifyParams struct {
	Username string
	Name     string
	Email    string
	Role     string
}

func HandleModifyUser(w http.ResponseWriter, r *http.Request, _ middlewares.Next) {
	userID := r.PathValue("id")
	decoder := json.NewDecoder(r.Body)
	params := modifyParams{}

	if err := decoder.Decode(&params); err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprint("Error in parsing json: ", err.Error()))
		return
	}

	if isEmpty := helpers.CheckEmptyStrings([]string{params.Email, params.Username, params.Role}); isEmpty {
		helpers.RespondWithError(w, 400, "email, username and role are required")
		return
	}

	if err := validators.CheckValidRole(params.Role); err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}

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
		helpers.RespondWithError(w, 400, "cann't modify account with this route")
		return
	}
	user := db.User{}
	// find user by id
	if res, err := ctx.Dbfn.FindOne("users", bson.D{{
		Key:   "_id",
		Value: objId,
	}}); err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	} else {
		if err = res.Decode(&user); err != nil {
			if strings.Contains(err.Error(), "no documents in result") {
				helpers.RespondWithError(w, 404, "user not found")
				return
			}
			helpers.RespondWithError(w, 400, err.Error())
			return
		}
	}

	// check if user is not admin
	if isAdmin := user.Role == "admin"; isAdmin {
		helpers.RespondWithError(w, 400, "cann't modify admin details")
		return
	}

	// update user
	if _, err := ctx.Dbfn.UpdateById("users", user.ID, bson.D{
		{Key: "name", Value: params.Name},
		{Key: "email", Value: params.Email},
		{Key: "role", Value: params.Role},
		{Key: "username", Value: params.Username},
	}); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			helpers.RespondWithError(w, 400, "email already taken")
			return
		}
		helpers.RespondWithError(w, 500, err.Error())
		return
	}

	helpers.RespondWithJSON(w, 200, "update successfully")
}
