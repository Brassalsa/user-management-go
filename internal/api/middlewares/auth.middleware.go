package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Brassalsa/user-management-go/internal"

	"github.com/Brassalsa/user-management-go/internal/db"
	"github.com/Brassalsa/user-management-go/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson"
)

type params struct {
	Token string `json:"token"`
}

type AuthCtx struct {
	Dbfn     *db.Database
	AuthUser *db.User
}

func VerifyTokenfunc(w http.ResponseWriter, r *http.Request, mCtx *AuthCtx, next *func()) {

	dbfn := mCtx.Dbfn
	decoder := json.NewDecoder(r.Body)
	tokenParam := params{}
	err := decoder.Decode(&tokenParam)

	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprint("error in token: ", err))
		return
	}

	authUser, err := internal.ValidateToken(tokenParam.Token)

	if err != nil {
		helpers.RespondWithError(w, 401, fmt.Sprint("Unauthorized: ", err))
		return
	}
	res, err := dbfn.FindOne("users", bson.D{{
		Key:   "username",
		Value: authUser.Username,
	}, {
		Key:   "_id",
		Value: authUser.Id,
	}, {
		Key:   "email",
		Value: authUser.Email,
	}})

	if err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}

	user := db.User{}

	res.Decode(&user)

	// call next
	(*next)()
}
