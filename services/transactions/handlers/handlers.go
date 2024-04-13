package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

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
			transactionJSON, err := json.Marshal(transaction)
			req, err := http.NewRequest("POST", "http://localhost:8080/commissions/calculate", bytes.NewBuffer(transactionJSON))
			if err != nil {
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
				fmt.Println(resp.Body)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			transaction.Amount += commission.Commission

			if err := repo.Create(transaction, db); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				fmt.Println(err)
				return
			}

			json.NewEncoder(w).Encode(models.TransactionResponse{Transaction: transaction, Ok: true})

		case "GET":
			id := r.PathValue("id")
			if id != "" {
				transaction := repo.Read(id, db)

				if transaction != nil {
					json.NewEncoder(w).Encode(models.TransactionResponse{Transaction: *transaction, Ok: true})
				} else {
					http.NotFound(w, r)
				}
			} else {
				transactions := repo.List(db)
				json.NewEncoder(w).Encode(models.ListResponse{Transaction: transactions, Ok: true})
			}
		case "UPDATE":
			id := r.PathValue("id")
			var transaction models.Transaction
			if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err := repo.Update(id, transaction, db); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		case "DELETE":
			id := r.PathValue("id")

			if err := repo.Delete(id, db); err != nil {
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
