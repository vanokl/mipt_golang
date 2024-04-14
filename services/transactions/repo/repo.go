package repo

import (
	"database/sql"
	"fmt"

	configs "github.com/vanokl/trxservice/config"
	"github.com/vanokl/trxservice/services/transactions/models"

	_ "github.com/lib/pq"
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
		fmt.Println("Exec err")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Ping err")
		return nil, err
	}

	return db, nil
}

func Create(transaction models.Transaction, db *sql.DB) error {

	_, err := db.Exec("INSERT INTO transactions (id, user_id, amount, currency, type, category, date, description) VALUES (DEFAULT, $1, $2, $3, $4, $5, $6, $7)", transaction.UserID, transaction.Amount, transaction.Currency, transaction.Type, transaction.Category, transaction.Date, transaction.Description)
	if err != nil {
		fmt.Println("Exec INSERT")
		return err
	}

	return nil
}

func Read(id string, db *sql.DB) *models.Transaction {
	var result models.Transaction

	row := db.QueryRow(`SELECT id, user_id, amount, currency, type, category, date, description FROM transactions WHERE id = $1`, id)
	if err := row.Scan(&result.ID, &result.UserID, &result.Amount, &result.Currency, &result.Type, &result.Category, &result.Date, &result.Description); err != nil {
		fmt.Println(err)
		return nil
	}

	return &result
}

func Update(id string, transaction models.Transaction, db *sql.DB) error {

	_, err := db.Exec(`UPDATE transactions 
					   SET user_id=$2, amount=$3, currency=$4, type=$5, category=$6, date=$7, description=$8
					   WHERE id=$1
					  `, id, transaction.UserID, transaction.Amount, transaction.Currency, transaction.Type, transaction.Category, transaction.Date, transaction.Description)
	if err != nil {
		fmt.Println("Exec UPDATE")
		return err
	}

	return nil
}

func Delete(id string, db *sql.DB) error {

	_, err := db.Exec(`DELETE FROM transactions WHERE id=$1`, id)
	if err != nil {
		fmt.Println("Exec DELETE")
		return err
	}

	return nil
}

func List(db *sql.DB) []models.Transaction {

	rows, err := db.Query(`SELECT * FROM transactions`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var results []models.Transaction

	for rows.Next() {
		var result models.Transaction
		if err := rows.Scan(&result.ID, &result.UserID, &result.Amount, &result.Currency, &result.Type, &result.Category, &result.Date, &result.Description); err != nil {
			fmt.Println(err)
			return nil
		}
		fmt.Println(result.ID)
		results = append(results, result)
	}

	return results
}
