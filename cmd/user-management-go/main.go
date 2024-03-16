package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Brassalsa/user-management-go/internal/api"
	"github.com/Brassalsa/user-management-go/internal/db"
)

func main() {
	// db connect
	dbFn := db.Database{
		Url: "mongodb://localhost:27017",
	}
	err := dbFn.Connect(context.TODO(), "user-management")

	if err != nil {
		log.Fatal(err)
	}

	// static files
	fs := http.FileServer(http.Dir("static/"))

	// router
	r := http.NewServeMux()
	// static files
	r.Handle("/static/", http.StripPrefix("/static/", fs))

	// testing endpoints
	r.HandleFunc("GET /healthz", api.HandlerReadiness)
	r.HandleFunc("GET /err", api.HandlerError)

	// api routes
	api.HandleV1Router(r, &dbFn)

	fmt.Println("server @ http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
