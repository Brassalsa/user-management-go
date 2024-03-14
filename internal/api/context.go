package api

import (
	"net/http"
)

type ApiContext[T any] struct {
	ctx T
}
type handlerFn[T any] func(w http.ResponseWriter, r *http.Request, ctx T)

func (ctx *ApiContext[T]) Provider(handler handlerFn[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, ctx.ctx)
	}
}
