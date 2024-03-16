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

type LoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// login user
func HandleLoginUser(w http.ResponseWriter, r *http.Request, dbfn *db.Database) {
	decoder := json.NewDecoder(r.Body)
	params := LoginParams{}
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprint("Error parsing json: ", err))
		return
	}

	isEmpty := helpers.CheckEmptyStrings([]string{params.Username, params.Password})

	if isEmpty {
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

	comparePass := internal.CompareHash(user.Password, params.Password)

	if !comparePass {
		helpers.RespondWithError(w, 400, "Wrong Credentials")
		return
	}
	token, err := internal.GenerateJWT(internal.AuthUser{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
	if err != nil {
		helpers.RespondWithError(w, 500, "Failed to Generate token")
		return
	}

	helpers.RespondWithJSON(w, 200, LoginResponse{
		Token: token,
	})
}

// rgister user
func HandleRegisterUser(w http.ResponseWriter, r *http.Request, dbfn *db.Database) {
	decoder := json.NewDecoder(r.Body)
	params := db.UserRegister{}
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprint("Error parsing json: ", err))
		return
	}

	isEmpty := helpers.CheckEmptyStrings([]string{params.Username, params.Password, params.Email})

	if isEmpty {
		helpers.RespondWithError(w, 400, "email, username, password are required")
		return
	}
	err = validators.CheckValidPassword(params.Password)

	if err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}

	params.Role = "user"
	err = internal.HashString(&params.Password)

	if err != nil {
		helpers.RespondWithError(w, 500, "Something went wrong")
		return
	}

	// inserting data
	_, err = dbfn.InsertIntoCollection("users", params)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			helpers.RespondWithError(w, 400, "username or email already exists")
			return
		}
		helpers.RespondWithError(w, 500, "Registering failed!")
		return
	}

	helpers.RespondWithJSON(w, 201, "User created successfully!")
}

// handle secure routes
type UserData struct {
	Id       primitive.ObjectID `json:"_id"`
	Username string             `json:"username"`
	Name     string             `json:"name"`
	Email    string             `json:"email"`
	Avatar   string             `json:"avatar"`
}

func HandleCheckCurrentUser(w http.ResponseWriter, r *http.Request, mCtx *middlewares.AuthCtx, next *func()) {
	authUser := mCtx.AuthUser

	helpers.RespondWithJSON(w, 200, UserData{
		Id:       authUser.ID,
		Username: authUser.Username,
		Name:     authUser.Name,
		Email:    authUser.Email,
		Avatar:   authUser.Avatar,
	})
}

// upload file
const MAX_UPLOAD_SIZE = 1024 * 1024 * 20

func HandleUploadAvatar(w http.ResponseWriter, r *http.Request, mCtx *middlewares.AuthCtx, next *func()) {

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		helpers.RespondWithError(w, 400, "The uploaded file is too big. Please choose an file that's less than 20MB in size")
		return
	}
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
	if isEmpty := helpers.CheckEmptyStrings([]string{mCtx.AuthUser.Avatar}); !isEmpty {
		if err := DeleteFile(mCtx.AuthUser.Avatar); err != nil {
			fmt.Println("err in deleting: ", err)
		}
	}
	mCtx.AuthUser.Avatar = url

	helpers.RespondWithJSON(w, 200, url)
}
