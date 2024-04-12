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
	CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL 
	);

	INSERT INTO userts VALUES (1, 'Иван', 'test@mail.ru', 'dfasdhfdgasdf');
	
	CREATE TYPE transaction_type AS ENUM ('income', 'expense', 'transfer');
		
	CREATE TABLE transactions (
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

	_, err := db.Exec("INSERT INTO transactions (id, user_id, amount, currency, type, category, date, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", transaction.ID, transaction.Amount, transaction.Currency, transaction.Type, transaction.Category, transaction.Date, transaction.Description)
	if err != nil {
		fmt.Println("Exec INSERT")
		return err
	}

	return nil
}

func Read(id string, db *sql.DB) *models.Transaction {
	var result models.Transaction
	db.QueryRow("SELECT * FROM transactions WHERE id = $1", id).Scan(&result)

	return &result
}

// Пример с использованием InMemoryDB
// var (
// 	errNotFound = errors.New("item not found")
// )

// type InMemoryDB struct {
// 	items map[string]models.Item
// 	mu    sync.RWMutex
// }

// func NewInMemoryDB() *InMemoryDB {
// 	return &InMemoryDB{
// 		items: make(map[string]models.Item),
// 	}
// }

// func (db *InMemoryDB) Create(item models.Item) {
// 	db.mu.Lock()
// 	defer db.mu.Unlock()

// 	db.items[item.ID] = item
// }

// func (db *InMemoryDB) Read(id string) (*models.Item, error) {
// 	db.mu.RLock()
// 	defer db.mu.RUnlock()

// 	item, exist := db.items[id]
// 	if !exist {
// 		return nil, errNotFound
// 	}

// 	return &item, nil
// }

// func (db *InMemoryDB) Update(id string, newValue string) error {
// 	db.mu.Lock()
// 	defer db.mu.Unlock()

// 	item, exist := db.items[id]
// 	if !exist {
// 		return errNotFound
// 	}

// 	item.Value = newValue
// 	db.items[id] = item

// 	return nil
// }

// func (db *InMemoryDB) Delete(id string) error {
// 	db.mu.Lock()
// 	defer db.mu.Unlock()

// 	_, exist := db.items[id]
// 	if !exist {
// 		return errNotFound
// 	}

// 	delete(db.items, id)

// 	return nil
// }

// func (db *InMemoryDB) List() []models.Item {
// 	db.mu.RLock()
// 	defer db.mu.RUnlock()

// 	items := make([]models.Item, 0, len(db.items))
// 	for _, item := range db.items {
// 		items = append(items, item)
// 	}

// 	return items
// }
