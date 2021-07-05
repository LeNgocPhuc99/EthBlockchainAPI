package main

import (
	"log"
	"net/http"

	"github.com/LeNgocPhuc99/blockchain-api/routes"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Server will start at http://localhost:8080/")

	route := mux.NewRouter()

	routes.AddRoutes(route)

	log.Fatal(http.ListenAndServe(":8080", route))
}
