package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Brassalsa/user-management-go/api"
)

func main() {
	r := http.DefaultServeMux
	
	
	// testing endpoints
	r.HandleFunc("GET /healthz" , api.HandlerReadiness)
	r.HandleFunc("GET /err" , api.HandlerError)

	// api routes
	api.HandleV1Router(r)
	
	fmt.Println("server @ http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))	
}