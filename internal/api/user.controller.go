package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Brassalsa/user-management-go/internal"
	"github.com/Brassalsa/user-management-go/internal/api/middlewares"
	"github.com/Brassalsa/user-management-go/internal/db"
	"github.com/Brassalsa/user-management-go/internal/validators"
	"github.com/Brassalsa/user-management-go/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type loginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// login user
func HandleLoginUser(w http.ResponseWriter, r *http.Request, dbfn *db.Database) {
	decoder := json.NewDecoder(r.Body)
	params := loginParams{}

	if err := decoder.Decode(&params); err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprint("Error parsing json: ", err))
		return
	}

	if isEmpty := helpers.CheckEmptyStrings([]string{params.Username, params.Password}); isEmpty {
		helpers.RespondWithError(w, 400, "username, password are required")
		return
	}

	res, err := dbfn.FindOne("users", bson.D{{
		Key:   "username",
		Value: params.Username,
	}})

	if err != nil {
		helpers.RespondWithError(w, 400, err.Error())
	}

	user := db.User{}
	res.Decode(&user)

	if comparePass := internal.CompareHash(user.Password, params.Password); !comparePass {
		helpers.RespondWithError(w, 400, "Wrong Credentials")
		return
	}
	token, err := internal.GenerateJWT(internal.AuthUser{
		Id:       user.ID,
		Username: user.Username,
	})
	if err != nil {
		helpers.RespondWithError(w, 500, "Failed to Generate token")
		return
	}

	helpers.RespondWithJSON(w, 200, loginResponse{
		Token: token,
	})
}

// register user
func HandleRegisterUser(w http.ResponseWriter, r *http.Request, dbfn *db.Database) {
	decoder := json.NewDecoder(r.Body)
	params := db.UserRegister{}
	if err := decoder.Decode(&params); err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprint("Error in parsing json: ", err.Error()))
		return
	}

	if isEmpty := helpers.CheckEmptyStrings([]string{params.Username, params.Password, params.Email}); isEmpty {
		helpers.RespondWithError(w, 400, "email, username, password are required")
		return
	}

	if err := validators.CheckValidPassword(params.Password); err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}

	params.Role = "user"

	if err := internal.HashString(&params.Password); err != nil {
		helpers.RespondWithError(w, 500, "Something went wrong")
		return
	}

	// inserting data

	if _, err := dbfn.InsertIntoCollection("users", params); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			helpers.RespondWithError(w, 400, "username or email already exists")
			return
		}
		helpers.RespondWithError(w, 500, fmt.Sprint("Registering failed: ", err.Error()))
		return
	}

	helpers.RespondWithJSON(w, 201, "User created successfully!")
}

// handle secure routes
type userData struct {
	Id       primitive.ObjectID `json:"_id"`
	Username string             `json:"username"`
	Name     string             `json:"name"`
	Email    string             `json:"email"`
	Avatar   string             `json:"avatar"`
}

func HandleCheckCurrentUser(w http.ResponseWriter, r *http.Request, mCtx *middlewares.AuthCtx, _ *func()) {
	authUser := mCtx.AuthUser

	helpers.RespondWithJSON(w, 200, userData{
		Id:       authUser.ID,
		Username: authUser.Username,
		Name:     authUser.Name,
		Email:    authUser.Email,
		Avatar:   authUser.Avatar,
	})
}

// upload file

func HandleUploadAvatar(w http.ResponseWriter, r *http.Request, mCtx *middlewares.AuthCtx, _ *func()) {

	url, err := UploadFile(r)

	if err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}
	user := mCtx.AuthUser
	// update avatar
	if _, err := mCtx.Dbfn.UpdateById("users", user.ID, bson.D{{
		Key:   "avatar",
		Value: url,
	}}); err != nil {
		helpers.RespondWithError(w, 500, err.Error())
		return
	}
	if isEmpty := helpers.CheckEmptyStrings([]string{user.Avatar}); !isEmpty {
		if err := DeleteFile(mCtx.AuthUser.Avatar); err != nil {
			fmt.Println("err in deleting: ", err)
		}
	}
	mCtx.AuthUser.Avatar = url

	helpers.RespondWithJSON(w, 200, url)
}

// update user's name and email
type updateParams struct {
	Name  string
	Email string
}

func HandleUpdateUser(w http.ResponseWriter, r *http.Request, mtCtx *middlewares.AuthCtx, _ *func()) {
	decoder := json.NewDecoder(r.Body)
	updateUser := updateParams{}
	if err := decoder.Decode(&updateUser); err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprint("Error in parsing json: ", err.Error()))
		return
	}

	if isEmpty := helpers.CheckEmptyStrings([]string{updateUser.Email}); isEmpty {
		helpers.RespondWithError(w, 400, "email is required")
		return
	}

	user := mtCtx.AuthUser

	if _, err := mtCtx.Dbfn.UpdateById("users", user.ID, bson.D{
		{Key: "name", Value: updateUser.Name},
		{Key: "email", Value: updateUser.Email},
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

// update password
type updatePassword struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func HandleUpdatePassword(w http.ResponseWriter, r *http.Request, mtCtx *middlewares.AuthCtx, _ *func()) {
	decoder := json.NewDecoder(r.Body)
	passwordParams := updatePassword{}
	if err := decoder.Decode(&passwordParams); err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprint("Error in parsing json: ", err.Error()))
		return
	}

	if isEmpty := helpers.CheckEmptyStrings([]string{passwordParams.OldPassword, passwordParams.NewPassword}); isEmpty {
		helpers.RespondWithError(w, 400, "oldPassword and newPassword are required")
		return
	}

	if err := validators.CheckValidPassword(passwordParams.NewPassword); err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}

	if checkPassword := internal.CompareHash(mtCtx.AuthUser.Password, passwordParams.OldPassword); !checkPassword {
		helpers.RespondWithError(w, 400, "wrong old password")
		return
	}

	if err := internal.HashString(&passwordParams.NewPassword); err != nil {
		helpers.RespondWithError(w, 500, fmt.Sprint("somehting went wrong: ", err.Error()))
		return
	}
	if _, err := mtCtx.Dbfn.UpdateById("users", mtCtx.AuthUser.ID, bson.D{
		{Key: "password", Value: passwordParams.NewPassword},
	}); err != nil {
		helpers.RespondWithError(w, 500, err.Error())
		return
	}

	helpers.RespondWithJSON(w, 200, "password updated")

}

// delete user account
type deleteUserParams struct {
	Username string
	Password string
}

func HandleDeleteUser(w http.ResponseWriter, r *http.Request, mtCtx *middlewares.AuthCtx, _ *func()) {
	decoder := json.NewDecoder(r.Body)
	params := deleteUserParams{}

	if err := decoder.Decode(&params); err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprint("Error in parsing json: ", err.Error()))
		return
	}

	if isEmpty := helpers.CheckEmptyStrings([]string{params.Username, params.Password}); isEmpty {
		helpers.RespondWithError(w, 400, "username and password are required")
		return
	}
	if err := validators.CheckValidPassword(params.Password); err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}

	user := mtCtx.AuthUser
	if checkPass := internal.CompareHash(user.Password, params.Password); !checkPass || user.Username != params.Username {
		helpers.RespondWithError(w, 400, "wrong credentials provided")
		return
	}

	if err := mtCtx.Dbfn.DeleteFromCollection("users", bson.D{{
		Key:   "_id",
		Value: user.ID,
	}}); err != nil {
		helpers.RespondWithError(w, 500, err.Error())
		return
	}

	helpers.RespondWithJSON(w, 200, "deleted successfully")

}
