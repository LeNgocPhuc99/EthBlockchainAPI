package routes

import (
	"fmt"

	"github.com/LeNgocPhuc99/blockchain-api/handlers"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
)

func AddRoutes(route *mux.Router) {
	client, err := ethclient.Dial("http://127.0.0.1:7545")

	if err != nil {
		fmt.Println(err)
	}

	route.Handle("/api/v1/eth/{module}", handlers.ClientHandler{Client: client})
}
