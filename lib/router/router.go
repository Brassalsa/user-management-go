package router

import (
	"net/http"
	"regexp"
)

type Route struct {
	Method  string
	Pattern string
	Handler http.Handler
}

type Router struct {
	routes []Route
}

type Handler func(w http.ResponseWriter, r *http.Request)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request){}


func (r *Router) AddRoute(method, path string, handler Handler) {
	r.routes = append(r.routes, Route{
		Method: method,
		Pattern: path,
		Handler: handler,
	})
}

func (r *Router) getHandler(method, path string) http.Handler {
	for _, route := range r.routes {
	   re := regexp.MustCompile(route.Pattern)
	   if route.Method == method && re.MatchString(path){
		  return route.Handler
	   }
	}
	return http.NotFoundHandler()
 }

 func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request){
	path := req.URL.Path
	method := req.Method
 
	handler := r.getHandler(method, path)
	handler.ServeHTTP(w, req)
 }

func (r *Router) GET(path string, handler Handler) {
	r.AddRoute("GET", path, handler)
 }
 
 func (r *Router) POST(path string, handler Handler) {
	r.AddRoute("POST", path, handler)
 }
 
 func (r *Router) PUT(path string, handler Handler) {
	r.AddRoute("PUT", path, handler)
 }
 
 func (r *Router) DELETE(path string, handler Handler) {
	r.AddRoute("DELETE", path, handler)
 }

 func (r *Router) Other(method, path string, handler Handler) {
	r.AddRoute(method, path, handler)
 }