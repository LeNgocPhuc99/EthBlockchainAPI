package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/LeNgocPhuc99/blockchain-api/models"
	"github.com/LeNgocPhuc99/blockchain-api/modules"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
)

type ClientHandler struct {
	Client *ethclient.Client
}

func (client ClientHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	module := vars["module"]

	// get query parameters
	address := r.URL.Query().Get("address")
	hash := r.URL.Query().Get("hash")

	rw.Header().Set("Content-Type", "application/json")

	switch module {
	case "latest-block":
		log.Println("get latest-block request")
		resBlock := modules.GetLatestBlock(*client.Client)
		json.NewEncoder(rw).Encode(resBlock)
	case "get-tx":
		log.Println("get transaction request")
		if hash == "" {
			json.NewEncoder(rw).Encode(&models.Error{
				Code:    http.StatusBadRequest,
				Message: "Bad request",
			})
			return
		}

		txHash := common.HexToHash(hash)
		resTx := modules.GetTxByHash(*client.Client, txHash)
		if resTx != nil {
			json.NewEncoder(rw).Encode(&models.Error{
				Code:    http.StatusNotFound,
				Message: "Tx Not Found !",
			})

			return
		}

		json.NewEncoder(rw).Encode(resTx)

	case "send-tx":
		log.Println("send-tx request")
		var transferReq models.TransferRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&transferReq)

		if err != nil {
			log.Println(err)
			json.NewEncoder(rw).Encode(models.Error{
				Code:    http.StatusBadRequest,
				Message: "Bad request",
			})

			return
		}

		// transfer eth
		hash, err := modules.TransferEth(*client.Client, transferReq.PrivateKey,
			transferReq.To, transferReq.Amount)

		if err != nil {
			fmt.Println(err)
			json.NewEncoder(rw).Encode(&models.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			})

			return
		}

		json.NewEncoder(rw).Encode(&models.HashResponse{
			Hash: hash,
		})

	case "get-balance":
		log.Println("get balance request")
		if address == "" {
			json.NewEncoder(rw).Encode(&models.Error{
				Code:    http.StatusBadRequest,
				Message: "Bad request",
			})

			return
		}

		balance, err := modules.GetAddressBalance(*client.Client, address)
		if err != nil {
			log.Println(err)
			json.NewEncoder(rw).Encode(&models.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			})
			return
		}

		json.NewEncoder(rw).Encode(&models.BalanceResponse{
			Address: address,
			Balance: balance,
			Symbol:  "Ether",
			Units:   "Wei",
		})
	}
}
