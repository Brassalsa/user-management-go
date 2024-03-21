package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Brassalsa/user-management-go/internal/api"
	"github.com/Brassalsa/user-management-go/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	portString := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URI")

	if portString == "" || dbUrl == "" {
		log.Fatal("PORT or DB_URI is not found in .env")
	}
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

	server := &http.Server{
		Handler: r,
		Addr:    ":" + portString}

	fmt.Println("server @ http://localhost:", portString)
	log.Fatal(server.ListenAndServe())
}
