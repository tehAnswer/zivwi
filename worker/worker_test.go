package worker_test

import (
	"encoding/json"
	"testing"

	nsq "github.com/bitly/go-nsq"
	"github.com/stretchr/testify/assert"
	app "github.com/tehAnswer/zivwi"
	worker "github.com/tehAnswer/zivwi/worker"
)

type FakeQueue struct{}

func (_ *FakeQueue) Publish(_ string, _ []byte) error {
	return nil
}

func TestDestAccountDontExist(t *testing.T) {
	database := app.NewDatabase()
	accountGateway := app.NewAccountGateway(database)
	transferGateway := app.NewTransferGateway(database)
	service := app.NewTransferService(
		accountGateway,
		transferGateway,
		&FakeQueue{})

	account1, _ := accountGateway.Create(app.Account{
		Balance: 10000,
	})

	transfer, _ := transferGateway.Create(app.Transfer{
		FromAccountId: account1.Id,
		ToAccountId:   "no_existe",
		Amount:        1,
		Message:       "fails.",
	})

	worker := worker.NewTransferWorker(accountGateway, transferGateway, service)
	payload, _ := json.Marshal(transfer)
	message := nsq.Message{Body: payload}
	if assert.Error(t, worker.HandleMessage(&message)) {
		transfer, _ = transferGateway.FindBy(transfer.Id)
		assert.Equal(t, "cancelled", transfer.Status)
		assert.Equal(t, "to_account_not_found", transfer.Error)
	}
}

func TestNotEnoughBalance(t *testing.T) {
	database := app.NewDatabase()
	accountGateway := app.NewAccountGateway(database)
	transferGateway := app.NewTransferGateway(database)
	service := app.NewTransferService(
		accountGateway,
		transferGateway,
		&FakeQueue{})

	account1, _ := accountGateway.Create(app.Account{
		Balance: 100000,
	})

	account2, _ := accountGateway.Create(app.Account{
		Balance: 0,
	})

	transfer, _ := transferGateway.Create(app.Transfer{
		FromAccountId: account1.Id,
		ToAccountId:   account2.Id,
		Amount:        9999999999999,
		Message:       "fails.",
	})

	worker := worker.NewTransferWorker(accountGateway, transferGateway, service)
	payload, _ := json.Marshal(transfer)
	message := nsq.Message{Body: payload}
	if assert.Error(t, worker.HandleMessage(&message)) {
		transfer, _ = transferGateway.FindBy(transfer.Id)
		assert.Equal(t, "cancelled", transfer.Status)
		assert.Equal(t, "not_enough_funds", transfer.Error)
	}
}

func TestSuccesfulTransfer(t *testing.T) {
	database := app.NewDatabase()
	accountGateway := app.NewAccountGateway(database)
	transferGateway := app.NewTransferGateway(database)
	service := app.NewTransferService(
		accountGateway,
		transferGateway,
		&FakeQueue{})

	account1, _ := accountGateway.Create(app.Account{
		Balance: 1,
	})

	account2, _ := accountGateway.Create(app.Account{
		Balance: 0,
	})

	transfer, _ := transferGateway.Create(app.Transfer{
		FromAccountId: account1.Id,
		ToAccountId:   account2.Id,
		Amount:        1,
		Message:       "not fails.",
	})

	worker := worker.NewTransferWorker(accountGateway, transferGateway, service)
	payload, _ := json.Marshal(transfer)
	message := nsq.Message{Body: payload}
	if assert.NoError(t, worker.HandleMessage(&message)) {
		transfer, _ = transferGateway.FindBy(transfer.Id)
		assert.Equal(t, "completed", transfer.Status)
		assert.Empty(t, transfer.Error)

		account1, _ := accountGateway.FindBy(account1.Id)
		account2, _ := accountGateway.FindBy(account2.Id)

		assert.Equal(t, uint64(0), account1.Balance)
		assert.Equal(t, uint64(1), account2.Balance)
	}
}
