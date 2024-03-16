package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Brassalsa/user-management-go/internal"

	"github.com/Brassalsa/user-management-go/internal/db"
	"github.com/Brassalsa/user-management-go/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson"
)

type AuthCtx struct {
	Dbfn     *db.Database
	AuthUser *db.User
}

func VerifyToken(w http.ResponseWriter, r *http.Request, mCtx *AuthCtx, next *func()) {
	dbfn := mCtx.Dbfn

	token, err := GetAuthToken(r.Header)
	if err != nil {
		helpers.RespondWithError(w, 400, "error in get token from header")
		return
	}

	authUser, err := internal.ValidateToken(token)

	if err != nil {
		helpers.RespondWithError(w, 401, "unauthorized token malformed")
		return
	}

	res, err := dbfn.FindOne("users", bson.D{
		{
			Key:   "_id",
			Value: authUser.Id,
		},
		{
			Key:   "username",
			Value: authUser.Username,
		}})

	if err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}
	user := db.User{}

	if err = res.Decode(&user); err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			helpers.RespondWithError(w, 401, "unauthorized")
			return
		}
		helpers.RespondWithError(w, 400, fmt.Sprint("something went wrong: ", err.Error()))
		return
	}

	if isEmpty := helpers.CheckEmptyStrings([]string{user.Username, user.Email}); isEmpty {
		helpers.RespondWithError(w, 401, "Unauthorized")
		return
	}
	mCtx.AuthUser = &user

	// call next
	(*next)()
}
