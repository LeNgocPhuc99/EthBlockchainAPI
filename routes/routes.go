package routes

import (
	"fmt"
	"net/http"

	"github.com/LeNgocPhuc99/blockchain-api/handlers"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	client, err := ethclient.Dial("http://127.0.0.1:7545")

	if err != nil {
		fmt.Println(err)
	}

	clientHandler := handlers.NewClientHandler(client)
	router := mux.NewRouter()

	router.HandleFunc("/api/latest-bock", clientHandler.HandleGetLastBlock).Methods(http.MethodGet)
	router.HandleFunc("/api/get-tx", clientHandler.HandleGetTransaction).Methods(http.MethodGet)
	router.HandleFunc("/api/tranfer-eth", clientHandler.HandleTransferEth).Methods(http.MethodPost)
	router.HandleFunc("/api/get-balance", clientHandler.HandleGetBalance).Methods(http.MethodGet)

	return router
}
