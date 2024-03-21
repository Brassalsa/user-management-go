package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getSecret() []byte {
	var secretStr = os.Getenv("JWT_SECRET")
	if secretStr == "" {
		fmt.Println("JWT_SECRET is not found in env, using 'mysecret' as secret string")
		return []byte("mysecret")
	}
	return []byte(secretStr)
}

type AuthUser struct {
	Id       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
}

// generate token
func GenerateJWT(authUser AuthUser) (string, error) {
	authUserJson, err := json.Marshal(authUser)
	if err != nil {
		err = fmt.Errorf("err in marshaling: %s", err.Error())
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = string(authUserJson)
	claims["exp"] = time.Now().Add(time.Minute * 60 * 24).Unix()

	tokenString, err := token.SignedString(getSecret())

	if err != nil {
		err = fmt.Errorf("something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (AuthUser, error) {
	authUser := AuthUser{}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getSecret(), nil
	})
	if err != nil {
		return authUser, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userJson := claims["user"].(string)

		err := json.Unmarshal([]byte(userJson), &authUser)
		if err != nil {
			return authUser, fmt.Errorf("error in unmarshaling: %v", err)
		}
		return authUser, nil
	} else {
		return authUser, err
	}
}
