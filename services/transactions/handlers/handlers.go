package handlers

import (
	"database/sql"
	"encoding/json"
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
				return
			}

			if err := repo.Create(transaction, db); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(models.TransactionResponse{Transaction: transaction, Ok: true})

		case "GET":
			id := r.URL.Query().Get("id")
			if id != "" {
				transaction := repo.Read(id, db)
				if transaction != nil {
					json.NewEncoder(w).Encode(models.TransactionResponse{Transaction: *transaction, Ok: true})
				} else {
					http.NotFound(w, r)
				}
			} else {
				json.NewEncoder(w).Encode(models.TransactionResponse{Transaction: models.Transaction{}, Ok: false})
			}
		}
	}
}
