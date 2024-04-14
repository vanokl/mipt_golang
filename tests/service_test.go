package repo

import (
	"testing"

	_ "github.com/lib/pq"
	configs "github.com/vanokl/trxservice/config"
	"github.com/vanokl/trxservice/services/transactions/models"
	"github.com/vanokl/trxservice/services/transactions/repo"
)

func TestCreateIntegration(t *testing.T) {

	config, err := configs.LoadConfig("../config/config.yaml")
	if err != nil {
		panic(err)
	}

	db, err := repo.InitDB(config)
	if err != nil {
		panic(err)
	}
	trx_before := repo.List(db)

	transaction := models.Transaction{
		ID:          "",
		UserID:      "1",
		Amount:      5.0,
		Currency:    "RUB",
		Type:        "income",
		Category:    "food",
		Date:        "2023-01-01",
		Description: "test",
	}
	repo.Create(transaction, db)

	trx_after := repo.List(db)

	if len(trx_before) == len(trx_after) {
		t.Errorf("no rows appended")
	}
}

func TestDeleteIntegration(t *testing.T) {

	config, err := configs.LoadConfig("../config/config.yaml")
	if err != nil {
		panic(err)
	}

	db, err := repo.InitDB(config)
	if err != nil {
		panic(err)
	}

	trx_before := repo.List(db)
	repo.Delete(trx_before[0].ID, db)
	trx_after := repo.List(db)

	if len(trx_before) == len(trx_after) {
		t.Errorf("no rows deleted")
	}
}

func TestReadIntegration(t *testing.T) {

	config, err := configs.LoadConfig("../config/config.yaml")
	if err != nil {
		panic(err)
	}

	db, err := repo.InitDB(config)
	if err != nil {
		panic(err)
	}

	transaction := models.Transaction{
		ID:          "",
		UserID:      "1",
		Amount:      5.0,
		Currency:    "RUB",
		Type:        "expense",
		Category:    "food",
		Date:        "2023-01-01",
		Description: "test",
	}
	repo.Create(transaction, db)
	list_trx := repo.List(db)
	id := list_trx[len(list_trx)-1].ID

	expected := models.Transaction{
		ID:          id,
		UserID:      "1",
		Amount:      5.0,
		Currency:    "RUB",
		Type:        "income",
		Category:    "food",
		Date:        "2023-01-01",
		Description: "test",
	}

	saved_trx := repo.Read(id, db)

	if saved_trx.Amount != expected.Amount {
		t.Errorf("amount in transactions not equal")
	}
}

func TestUpdateIntegration(t *testing.T) {

	config, err := configs.LoadConfig("../config/config.yaml")
	if err != nil {
		panic(err)
	}

	db, err := repo.InitDB(config)
	if err != nil {
		panic(err)
	}

	transaction := models.Transaction{
		ID:          "",
		UserID:      "1",
		Amount:      5.0,
		Currency:    "RUB",
		Type:        "expense",
		Category:    "food",
		Date:        "2023-01-01",
		Description: "test",
	}
	repo.Create(transaction, db)
	list_trx := repo.List(db)
	id := list_trx[len(list_trx)-1].ID

	updated := models.Transaction{
		ID:          id,
		UserID:      "1",
		Amount:      10.0,
		Currency:    "RUB",
		Type:        "income",
		Category:    "food",
		Date:        "2023-01-01",
		Description: "test",
	}
	repo.Update(id, updated, db)

	expected := repo.Read(id, db)

	if expected.Amount != 10 {
		t.Errorf("amount in transactions not updated")
	}
}
