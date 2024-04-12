package main

import (
	"log"
	"net/http"

	configs "github.com/vanokl/trxservice/config"
	"github.com/vanokl/trxservice/services/transactions/handlers"
	"github.com/vanokl/trxservice/services/transactions/repo"
)

func main() {

	config, err := configs.LoadConfig("config/config.yaml")
	if err != nil {
		panic(err)
	}

	db, err := repo.InitDB(config)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/transactions", handlers.Transaction(db))
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
