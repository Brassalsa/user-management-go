package middlewares

import (
	"errors"
	"net/http"
	"strings"
)

type HFunc[T any] func(w http.ResponseWriter, r *http.Request, ctx T, next *func())

func GroupMiddlewares[T any](ctx T, handleFuncs []HFunc[T]) http.HandlerFunc {
	goNext := 0
	next := func() {
		goNext += 1
	}

	return func(w http.ResponseWriter, r *http.Request) {
		goNext = 0
		for ind, val := range handleFuncs {
			if goNext != ind || goNext == len(handleFuncs) {
				return
			}
			val(w, r, ctx, &next)
		}
	}
}

// get token from header:
// Authorization   token <your-token>
func GetAPiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no auth info found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}

	if vals[0] != "token" {
		return "", errors.New("malformed first part of auth header")
	}

	return vals[1], nil
}
