package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Brassalsa/user-management-go/api"
	"github.com/Brassalsa/user-management-go/internal/db"
)

func main() {

	r := http.DefaultServeMux
	db := db.Database{
		Url:         "mongodb://localhost:27017",
		Collections: []string{"users"},
	}
	db.Connect(context.TODO(), "user-management")
	// testing endpoints
	r.HandleFunc("GET /healthz", api.HandlerReadiness)
	r.HandleFunc("GET /err", api.HandlerError)

	// api routes
	api.HandleV1Router(r)

	fmt.Println("server @ http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
