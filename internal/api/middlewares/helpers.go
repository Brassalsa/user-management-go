package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type Next *func()

type HFunc func(w http.ResponseWriter, r *http.Request, next Next)

type ctxKey string

const CtxKey ctxKey = "ctx"

func GroupMiddlewares[T any](ctx T, handleFuncs []HFunc) http.HandlerFunc {
	goNext := 0
	next := func() {
		goNext += 1
	}

	return func(w http.ResponseWriter, r *http.Request) {
		goNext = 0
		mctx := context.WithValue(r.Context(), CtxKey, ctx)
		for ind, val := range handleFuncs {
			if goNext != ind || goNext == len(handleFuncs) {
				return
			}
			val(w, r.WithContext(mctx), &next)
		}
	}
}

// get token from header:
// Authorization   Bearer <your-token>
func GetAuthToken(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no auth info found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}

	if vals[0] != "Bearer" {
		return "", errors.New("malformed first part of auth header")
	}

	return vals[1], nil
}
