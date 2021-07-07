package main

import (
	"log"
	"net/http"

	"github.com/LeNgocPhuc99/blockchain-api/routes"
)

func main() {
	log.Println("Server will start at http://localhost:8080/")

	r := routes.Router()

	log.Fatal(http.ListenAndServe(":8080", r))
}
