package repo

import (
	"database/sql"
	"fmt"

	"log"

	_ "github.com/lib/pq"
	configs "github.com/vanokl/trxservice/config"
	"github.com/vanokl/trxservice/services/transactions/models"
)

func InitDB(config *configs.Config) (*sql.DB, error) {
	dbConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.DBName, config.Database.SSLMode)
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		fmt.Println("Open err")
		return nil, err
	}

	createDB := `
	CREATE TABLE IF NOT EXISTS users  (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL 
	);

	--INSERT INTO users VALUES (1, 'Иван', 'test@mail.ru', 'dfasdhfdgasdf');
	
	--CREATE OR REPLACE TYPE  transaction_type  AS ENUM ('income', 'expense', 'transfer');
		
	CREATE TABLE IF NOT EXISTS transactions   (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id),
		amount NUMERIC(10, 2) NOT NULL CHECK (amount >= 0),
		currency VARCHAR(255) NOT NULL,
		type transaction_type NOT NULL,
		category VARCHAR(255) NOT NULL,
		date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
		description TEXT
	);
	`

	_, err = db.Exec(createDB)
	if err != nil {
		log.Fatal("Exec err")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Ping err")
		return nil, err
	}

	return db, nil
}

func Create(transaction models.Transaction, db *sql.DB) error {
	log.Println("Exec INSERT")
	_, err := db.Exec("INSERT INTO transactions (id, user_id, amount, currency, type, category, date, description) VALUES (DEFAULT, $1, $2, $3, $4, $5, $6, $7)", transaction.UserID, transaction.Amount, transaction.Currency, transaction.Type, transaction.Category, transaction.Date, transaction.Description)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func Read(id string, db *sql.DB) *models.Transaction {
	var result models.Transaction
	log.Println("Exec SELECT")
	row := db.QueryRow(`SELECT id, user_id, amount, currency, type, category, date, description FROM transactions WHERE id = $1`, id)
	if err := row.Scan(&result.ID, &result.UserID, &result.Amount, &result.Currency, &result.Type, &result.Category, &result.Date, &result.Description); err != nil {
		log.Fatal(err)
		return nil
	}

	return &result
}

func Update(id string, transaction models.Transaction, db *sql.DB) error {
	log.Println("Exec UPDATE")
	_, err := db.Exec(`UPDATE transactions 
					   SET user_id=$2, amount=$3, currency=$4, type=$5, category=$6, date=$7, description=$8
					   WHERE id=$1
					  `, id, transaction.UserID, transaction.Amount, transaction.Currency, transaction.Type, transaction.Category, transaction.Date, transaction.Description)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func Delete(id string, db *sql.DB) error {
	log.Println("Exec DELETE")
	_, err := db.Exec(`DELETE FROM transactions WHERE id=$1`, id)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func List(db *sql.DB) []models.Transaction {
	log.Println("Exec LIST")
	rows, err := db.Query(`SELECT * FROM transactions`)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer rows.Close()

	var results []models.Transaction

	for rows.Next() {
		var result models.Transaction
		if err := rows.Scan(&result.ID, &result.UserID, &result.Amount, &result.Currency, &result.Type, &result.Category, &result.Date, &result.Description); err != nil {
			log.Fatal(err)
			return nil
		}
		fmt.Println(result.ID)
		results = append(results, result)
	}

	return results
}
