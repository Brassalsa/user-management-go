package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// testing endpoints
	http.HandleFunc("GET /healthz" ,handlerReadiness)
	http.HandleFunc("GET /err" ,handlerError)

	fmt.Println("server @ http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}