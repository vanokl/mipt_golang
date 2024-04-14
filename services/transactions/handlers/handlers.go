package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/vanokl/trxservice/services/transactions/models"
	"github.com/vanokl/trxservice/services/transactions/repo"
)

func Transaction(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":

			var transaction models.Transaction
			if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				fmt.Println(err)
				return
			}
			log.Println("Calculate commission")
			transactionJSON, err := json.Marshal(transaction)
			if err != nil {
				log.Fatal(err)
				return
			}

			req, err := http.NewRequest("POST", "http://localhost:8080/commissions/calculate", bytes.NewBuffer(transactionJSON))
			if err != nil {
				log.Fatal(err)
				return
			}

			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			var commission models.Commission
			if err := json.NewDecoder(resp.Body).Decode(&commission); err != nil {
				log.Fatal(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			transaction.Amount += commission.Commission

			log.Println("save transaction")
			if err := repo.Create(transaction, db); err != nil {
				log.Fatal(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(models.TransactionResponse{Transaction: transaction, Ok: true})

		case "GET":

			id := r.PathValue("id")
			currency := r.URL.Query().Get("currency")
			if currency != "" {
				log.Println("convert currency")
				transaction := repo.Read(id, db)

				url := fmt.Sprintf("http://localhost:8080/convert?amount=%f&from=%s&to=%s", transaction.Amount, transaction.Currency, currency)
				log.Println(url)
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					log.Fatal(err)
					return
				}

				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					log.Fatal(err)
					return
				}
				defer resp.Body.Close()

				var converted models.CurrencyAnswer
				if err := json.NewDecoder(resp.Body).Decode(&converted); err != nil {
					log.Fatal(err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				converted_transaction := models.TransactionConverted{
					ID:                transaction.ID,
					UserID:            transaction.UserID,
					Amount:            transaction.Amount,
					Currency:          transaction.Currency,
					ConvertedAmount:   converted.ConvertedAmount,
					ConvertedCurrency: converted.To,
					Type:              transaction.Type,
					Category:          transaction.Category,
					Date:              transaction.Date,
					Description:       transaction.Description,
				}

				json.NewEncoder(w).Encode(converted_transaction)

			} else {
				if id != "" {
					log.Println("get transaction")
					transaction := repo.Read(id, db)

					if transaction != nil {
						json.NewEncoder(w).Encode(models.TransactionResponse{Transaction: *transaction, Ok: true})
					} else {
						log.Fatal("transaction not found")
						http.NotFound(w, r)
					}
				} else {
					log.Println("list all transactions")
					transactions := repo.List(db)
					json.NewEncoder(w).Encode(models.ListResponse{Transaction: transactions, Ok: true})
				}
			}

		case "UPDATE":
			log.Println("update transaction")
			id := r.PathValue("id")
			var transaction models.Transaction
			if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
				log.Fatal(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err := repo.Update(id, transaction, db); err != nil {
				log.Fatal(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		case "DELETE":

			log.Println("delete transaction")
			id := r.PathValue("id")

			if err := repo.Delete(id, db); err != nil {
				log.Fatal(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

func CommissionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			var transaction models.Transaction
			if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var commissionRate float64
			switch transaction.Type {
			case "income":
				switch transaction.Currency {
				case "USD":
					commissionRate = 0.02
				case "EUR":
					commissionRate = 0.04
				case "RUB":
					commissionRate = 0.05
				}
			}
			commission := models.Commission{
				TransactionID:   transaction.ID,
				Amount:          transaction.Amount,
				Currency:        transaction.Currency,
				Type:            transaction.Type,
				Commission:      commissionRate * transaction.Amount,
				CalculationDate: transaction.Date,
				Description:     "Commission for transaction",
			}

			json.NewEncoder(w).Encode(commission)

		}
	}
}

func CurrencyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Println("get currency from https://api.freecurrencyapi.com")
		switch r.Method {
		case "GET":
			from := r.URL.Query().Get("from")
			to := r.URL.Query().Get("to")
			amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
			if err != nil {
				log.Println(err)
			}
			url := fmt.Sprintf("https://api.freecurrencyapi.com/v1/latest?base_currency=%s&currencies=RUB,USD,EUR&apikey=fca_live_g25GMMWQ9VcXztgooX3yIietVivVGVJgmRDGcTMp", from)

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Println(err)
				return
			}

			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				return
			}
			defer resp.Body.Close()

			var currency models.CurrencyStruct
			if err := json.NewDecoder(resp.Body).Decode(&currency); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var rate float64
			switch to {
			case "USD":

				rate = currency.DataCurrency.USD
			case "EUR":
				rate = currency.DataCurrency.EUR
			case "RUB":
				rate = currency.DataCurrency.RUB

			}

			answer := models.CurrencyAnswer{
				From:            from,
				To:              to,
				OriginalAmount:  amount,
				ConvertedAmount: amount * rate,
				Rate:            rate,
			}

			json.NewEncoder(w).Encode(answer)
		}
	}
}
