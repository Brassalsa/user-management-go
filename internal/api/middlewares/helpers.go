package middlewares

import "net/http"

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
